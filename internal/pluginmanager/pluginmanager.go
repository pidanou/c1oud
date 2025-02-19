package pluginmanager

import (
	"errors"

	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/pkg/plugin"
)

type PluginManager struct {
	PluginRepository repositories.PluginRepository
}

func NewPluginManager(repo repositories.PluginRepository) *PluginManager {
	return &PluginManager{PluginRepository: repo}
}

func (p *PluginManager) InstallPlugin(plug *plugin.Plugin) (*plugin.Plugin, error) {
	if plug.URI == "" {
		return nil, errors.New("Missing URI")
	}

	if plug.Source != plugin.VCS && plug.Source != plugin.Local && plug.Source != plugin.HTTP {
		return nil, errors.New("Source not supported")
	}

	var err error
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
	return p.PluginRepository.DeletePlugin(name)
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

func (p *PluginManager) GetAccount(id int) (*plugin.Account, error) {
	return p.PluginRepository.GetAccount(id)
}

func (p *PluginManager) AddAccount(acc *plugin.Account) (*plugin.Account, error) {
	return p.PluginRepository.AddAccount(acc)
}

func (p *PluginManager) DeleteAccount(id int) error {
	return p.PluginRepository.DeleteAccount(id)
}
