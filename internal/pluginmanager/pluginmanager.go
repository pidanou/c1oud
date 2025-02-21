package pluginmanager

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/internal/types"
	"github.com/pidanou/c1-core/pkg/plugin"
)

type PluginManager struct {
	PluginRepository repositories.PluginRepository
}

func NewPluginManager(repo repositories.PluginRepository) *PluginManager {
	return &PluginManager{PluginRepository: repo}
}

func (p *PluginManager) GetPlugin(name string) (*plugin.Plugin, error) {
	return p.PluginRepository.GetPlugin(name)
}

func (p *PluginManager) InstallPlugin(pluginForm *types.PluginForm) (*plugin.Plugin, error) {
	plug := &plugin.Plugin{}

	err := json.Unmarshal([]byte(pluginForm.Config), plug)
	if err != nil {
		client := &http.Client{}

		req, err := http.NewRequest("GET", pluginForm.Config, nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println(err)
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		err = json.Unmarshal(body, plug)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if plug.URI == "" {
		return nil, errors.New("Missing URI")
	}

	if plug.Source != plugin.VCS && plug.Source != plugin.Local && plug.Source != plugin.HTTP {
		return nil, errors.New("Source not supported")
	}

	if pluginForm.NameOverride != "" {
		plug.Name = pluginForm.NameOverride
	}

	switch plug.Source {
	case plugin.VCS:
		err = downloadFromVCS(plug)
	case plugin.HTTP:
		err = downloadFromHTTP(plug)
	case plugin.Local:
		err = downloadFromLocal(plug)
	}

	if err != nil {
		return nil, err
	}
	// WIP Add cleanup if AddPlugin fails
	return p.PluginRepository.AddPlugin(plug)
}

func (p *PluginManager) DeletePlugin(name string) error {
	err := p.PluginRepository.DeletePlugin(name)
	if err != nil {
		log.Println(err)
		return errors.New("Unable to remove plugin from database")
	}
	err = DeletePlugin(name)
	if err != nil {
		return errors.New("Unable to delete plugin")
	}
	return nil
}

func (p *PluginManager) EditPlugin(plug *plugin.Plugin) (*plugin.Plugin, error) {
	return p.PluginRepository.EditPlugin(plug)
}

func (p *PluginManager) UpdatePlugin(name string) error {
	plug, err := p.PluginRepository.GetPlugin(name)
	if err != nil {
		return err
	}
	if plug.URI == "" {
		return errors.New("Missing URI")
	}

	if plug.Source != plugin.VCS && plug.Source != plugin.Local && plug.Source != plugin.HTTP {
		return errors.New("Source not supported")
	}

	switch plug.Source {
	case plugin.VCS:
		return updateFromVCS(plug)
	case plugin.HTTP:
		return updateFromHTTP(plug)
	case plugin.Local:
		return updateFromLocal(plug)
	}

	return nil
}

func (p *PluginManager) ListPlugins() ([]plugin.Plugin, error) {
	return p.PluginRepository.ListPlugins()
}

func (p *PluginManager) ListAccounts() ([]plugin.Account, error) {
	return p.PluginRepository.ListAccounts()
}

func (p *PluginManager) GetAccount(id int32) (*plugin.Account, error) {
	return p.PluginRepository.GetAccount(id)
}

func (p *PluginManager) AddAccount(acc *plugin.Account) (*plugin.Account, error) {
	return p.PluginRepository.AddAccount(acc)
}

func (p *PluginManager) EditAccount(account *plugin.Account) (*plugin.Account, error) {
	return p.PluginRepository.EditAccount(account)
}

func (p *PluginManager) DeleteAccount(id int32) error {
	return p.PluginRepository.DeleteAccount(id)
}

func (p *PluginManager) ListData(limit, offset int, filters map[string]string) ([]plugin.Data, error) {
	return p.PluginRepository.ListData(limit, offset, filters)
}

func (p *PluginManager) GetData(id int32) (*plugin.Data, error) {
	return p.PluginRepository.GetData(id)
}

func (p *PluginManager) EditData(data *plugin.Data) (*plugin.Data, error) {
	return p.PluginRepository.EditData(data)
}
