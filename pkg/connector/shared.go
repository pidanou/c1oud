package connector

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/pidanou/c1-core/pkg/connector/proto"
	"google.golang.org/grpc"
)

// 1. Client calls the implementation Sync(string) [see executor.go "sync" method]
// 2. rpc call: Sync(SyncRequest) returns SyncResponse on the GRPC server  [see plugin.proto]
// 3. Server calls the implementation of Sync [the plugin Sync method]
// 4. The server calls back and sends the response to another server that handles storage of the response.
//    This allows pagination for plugins.
// 5. Returns errors if any to the client

// Plugins need to implement this interface
type ConnectorInterface interface {
	Sync(options string, c CallbackHandler) error
}

type CallbackHandler interface {
	Callback(*proto.SyncResponse) (*proto.Empty, error)
}

// GRPC Client: started by go-plugin to call RPC methods
type GRPCConnectorClient struct {
	broker *plugin.GRPCBroker
	client proto.ConnectorClient
}

// Implementation of the Connector interface: the client calls this function that is sent to the GRPC server (the plugin).
func (g *GRPCConnectorClient) Sync(options string, c CallbackHandler) error {
	callbackHandlerServer := &GRPCCallbackHandlerServer{Impl: c}
	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterCallbackHandlerServer(s, callbackHandlerServer)

		return s
	}

	brokerID := g.broker.NextId()
	go g.broker.AcceptAndServe(brokerID, serverFunc)
	_, err := g.client.Sync(context.Background(), &proto.SyncRequest{
		Options:               options,
		CallbackHandlerServer: brokerID,
	})
	return err
}

// GRPC Server: started by go-plugin to listen to client RPC calls for the plugin
type ConnectorGRPCServer struct {
	proto.UnimplementedConnectorServer
	Impl   ConnectorInterface
	broker *plugin.GRPCBroker
}

// Implementation of the Connector interface: the servers calls the plugin implementation of the function
func (s *ConnectorGRPCServer) Sync(ctx context.Context,
	req *proto.SyncRequest) (*proto.Empty, error) {
	conn, err := s.broker.Dial(req.CallbackHandlerServer)
	if err != nil {
		return &proto.Empty{}, err
	}
	defer conn.Close()
	c := &GRPCCallbackHandlerClient{proto.NewCallbackHandlerClient(conn)}
	return &proto.Empty{}, s.Impl.Sync(req.Options, c)
}

// Connector plugin over GRPC
type ConnectorGRPCPlugin struct {
	plugin.Plugin
	Impl ConnectorInterface
}

// Build a GRPC Server for the ConnectorGRPC plugin. This server will run the plugin
func (p *ConnectorGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterConnectorServer(s, &ConnectorGRPCServer{Impl: p.Impl, broker: broker})
	return nil
}

// Build a GRPC Client for the ConnectorGRPC plugin
func (ConnectorGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCConnectorClient{client: proto.NewConnectorClient(c), broker: broker}, nil
}

// GRPC client/server to handle the data response
type GRPCCallbackHandlerClient struct{ client proto.CallbackHandlerClient }

func (m *GRPCCallbackHandlerClient) Callback(res *proto.SyncResponse) (*proto.Empty, error) {
	_, err := m.client.Callback(context.Background(), res)
	if err != nil {
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, nil
}

type GRPCCallbackHandlerServer struct {
	proto.UnimplementedCallbackHandlerServer
	Impl CallbackHandler
}

func (m *GRPCCallbackHandlerServer) Callback(ctx context.Context, req *proto.SyncResponse) (*proto.Empty, error) {
	return m.Impl.Callback(req)
}
