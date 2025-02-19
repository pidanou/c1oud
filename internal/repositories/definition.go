package repositories

import "github.com/pidanou/c1-core/pkg/plugin"

type PluginRepository interface {
	GetPlugin(name string) (*plugin.Plugin, error)
	AddPlugin(*plugin.Plugin) (*plugin.Plugin, error)
	DeletePlugin(name string) error
	GetAccount(id int) (*plugin.Account, error)
	AddAccount(*plugin.Account) (*plugin.Account, error)
	DeleteAccount(id int) error
	AddData([]plugin.Data) error
}
