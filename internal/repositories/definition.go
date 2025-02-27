package repositories

import (
	"github.com/pidanou/c1-core/internal/types"
	"github.com/pidanou/c1-core/pkg/connector"
)

type ConnectorRepository interface {
	ListActiveConnectors() (res []connector.Connector, count int, err error)
	ListAllConnectors() (res []string, err error)
	GetConnector(name string) (*connector.Connector, error)
	AddConnector(*connector.Connector) (*connector.Connector, error)
	EditConnector(*connector.Connector) (*connector.Connector, error)
	DeleteConnector(name string) error
	ListAccounts() (res []connector.Account, count int, err error)
	GetAccount(id int32) (*connector.Account, error)
	AddAccount(*connector.Account) (*connector.Account, error)
	EditAccount(*connector.Account) (*connector.Account, error)
	DeleteAccount(id int32) error
	ListData(filters *types.Filter) (res []connector.Data, count int, err error)
	AddData([]connector.Data) error
	GetData(id int32) (*connector.Data, error)
	EditData(*connector.Data) (*connector.Data, error)
}
