package repositories

import "github.com/pidanou/c1-core/pkg/plugin"

type PluginRepository interface {
	ListPlugins() ([]plugin.Plugin, error)
	GetPlugin(name string) (*plugin.Plugin, error)
	AddPlugin(*plugin.Plugin) (*plugin.Plugin, error)
	EditPlugin(*plugin.Plugin) (*plugin.Plugin, error)
	DeletePlugin(name string) error
	ListAccounts() ([]plugin.Account, error)
	GetAccount(id int32) (*plugin.Account, error)
	AddAccount(*plugin.Account) (*plugin.Account, error)
	EditAccount(*plugin.Account) (*plugin.Account, error)
	DeleteAccount(id int32) error
	ListData(limit, offset int, filters map[string]string) ([]plugin.Data, error)
	AddData([]plugin.Data) error
	GetData(id int32) (*plugin.Data, error)
	EditData(*plugin.Data) (*plugin.Data, error)
}
