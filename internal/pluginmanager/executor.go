package pluginmanager

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
	plug "github.com/pidanou/c1-core/pkg/plugin"
	"github.com/pidanou/c1-core/pkg/plugin/proto"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"connector": &plug.ConnectorGRPCPlugin{},
}

func (p *PluginManager) Execute(accountIDs []int) error {
	results := make(chan error, len(accountIDs))
	var wg sync.WaitGroup

	for _, accountID := range accountIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result := p.sync(id)
			results <- result
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

func (p *PluginManager) sync(accountID int) error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	acc, err := p.GetAccount(accountID)
	if err != nil {
		return err
	}
	pl, _ := p.PluginRepository.GetPlugin(acc.Plugin)
	cmd := exec.Command("sh", "-c", pl.Command)
	cmd.Dir = path.Join(constants.Envs["C1_DIR"], pl.Name)
	client := plugin.NewClient(&plugin.ClientConfig{
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
		log.Println(err)
	}

	raw, err := gRPCClient.Dispense("connector")
	if err != nil {
		log.Println(err)
	}

	connector := raw.(plug.Connector)
	err = connector.Sync(acc.Options, &callbackHandler{Plugin: pl.Name, DataRepository: p.PluginRepository})
	if err != nil {
		return err
	}

	return nil
}

type callbackHandler struct {
	Plugin         string
	DataRepository repositories.PluginRepository
}

func (c *callbackHandler) Callback(res *proto.SyncResponse) (*proto.Empty, error) {
	data := []plug.Data{}
	for _, obj := range res.Response {
		d := plug.Data{RemoteID: obj.RemoteId, Plugin: c.Plugin, ResourceName: obj.ResourceName, URI: obj.Uri, Metadata: obj.Metadata}
		data = append(data, d)
	}
	err := c.DataRepository.AddData(data)
	if err != nil {
		log.Println("err", err)
	}
	return &proto.Empty{}, err
}
