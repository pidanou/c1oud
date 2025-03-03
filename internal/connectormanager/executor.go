package connectormanager

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/pkg/connector"
	"github.com/pidanou/c1-core/pkg/connector/proto"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"connector": &connector.ConnectorGRPCPlugin{},
}

func (p *ConnectorManager) Execute(accountIDs []int32) error {
	results := make(chan error, len(accountIDs))
	var wg sync.WaitGroup

	for _, accountID := range accountIDs {
		wg.Add(1)
		go func(id int32) {
			defer wg.Done()
			account, err := p.ConnectorRepository.GetAccount(accountID)
			if err != nil {
				return
			}
			metadata, err := p.sync(id)
			results <- err
			syncInfo := &connector.SyncInfo{
				Connector: account.Connector,
				Metadata:  metadata,
				Success:   err == nil,
			}
			_ = p.ConnectorRepository.AddSyncInfo(syncInfo)
		}(accountID)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	errs := []error{}
	for result := range results {
		if result == nil {
			continue
		}
		errs = append(errs, result)
	}

	if len(errs) > 0 {
		log.Println(errs)
		return errors.New("Execution has met some errors")
	}

	return nil
}

func (p *ConnectorManager) sync(accountID int32) (string, error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "connector",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	acc, err := p.GetAccount(accountID)
	if err != nil {
		return "An error occured before synchronization", err
	}
	pl, _ := p.ConnectorRepository.GetConnector(acc.Connector)
	cmd := exec.Command("sh", "-c", pl.Command)
	cmd.Dir = path.Join(constants.Envs["C1_DIR"], pl.Name)
	client := plugin.NewClient(&plugin.ClientConfig{
		UnixSocketConfig: &plugin.UnixSocketConfig{
			TempDir: "/Users/pidanou/.c1/tmp",
		},
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
	})
	defer client.Kill()

	gRPCClient, err := client.Client()
	if err != nil {
		log.Println("cannot get client ", err)
	}
	if gRPCClient == nil {
		return "An error occured before synchronization", errors.New("no connector client")
	}

	raw, err := gRPCClient.Dispense("connector")
	if err != nil {
		log.Println(err)
		return "An error occured before synchronization", err
	}

	conn := raw.(connector.ConnectorInterface)
	metadata, err := conn.Sync(acc.Options, &callbackHandler{AccountID: acc.ID, Connector: pl.Name, DataRepository: p.ConnectorRepository})
	return metadata, err
}

type callbackHandler struct {
	AccountID      int32
	Connector      string
	DataRepository repositories.ConnectorRepository
}

func (c *callbackHandler) Callback(res *proto.SyncResponse) error {
	data := []connector.Data{}
	for _, obj := range res.Response {
		d := connector.Data{AccountID: c.AccountID, RemoteID: obj.RemoteId, Connector: c.Connector, ResourceName: obj.ResourceName, URI: obj.Uri, Metadata: obj.Metadata}
		data = append(data, d)
	}
	err := c.DataRepository.AddData(data)
	if err != nil {
		log.Println("err", err)
	}
	return err
}
