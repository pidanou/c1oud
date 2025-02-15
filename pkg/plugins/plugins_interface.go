package plugins

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Response struct {
}

// Greeter is the interface that we're exposing as a plugin.
type Connector interface {
	Sync(options interface{}) Response
}

// Here is an implementation that talks over RPC
type ConnectorRPC struct{ client *rpc.Client }

func (g *ConnectorRPC) Sync(name string) Response {
	var resp Response
	err := g.client.Call("Plugin.Greet", name, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type GreeterRPCServer struct {
	// This is the real implementation
	Impl Connector
}

func (s *GreeterRPCServer) Sync(options interface{}, resp Response) error {
	resp = s.Impl.Sync(options)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type ConnectorPlugin struct {
	// Impl Injection
	Impl Connector
}

func (p *ConnectorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl}, nil
}

func (ConnectorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ConnectorRPC{client: c}, nil
}
