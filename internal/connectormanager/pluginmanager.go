package connectormanager

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/internal/types"
	"github.com/pidanou/c1-core/pkg/connector"
)

type ConnectorManager struct {
	ConnectorRepository repositories.ConnectorRepository
}

func NewConnectorManager(repo repositories.ConnectorRepository) *ConnectorManager {
	return &ConnectorManager{ConnectorRepository: repo}
}

func (p *ConnectorManager) GetConnector(name string) (*connector.Connector, error) {
	return p.ConnectorRepository.GetConnector(name)
}

func (p *ConnectorManager) InstallConnector(connectorForm *types.ConnectorForm) (*connector.Connector, error) {
	conn := &connector.Connector{}

	err := json.Unmarshal([]byte(connectorForm.Config), conn)
	if err != nil {
		client := &http.Client{}

		req, err := http.NewRequest("GET", connectorForm.Config, nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		resp, err := client.Do(req)
		if err == nil {
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

			err = json.Unmarshal(body, conn)
			if err != nil {
				return nil, err
			}
		} else {
			// Try from local
			file, err := os.Open(connectorForm.Config)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			bytes, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(bytes, conn)
			if err != nil {
				return nil, err
			}
		}
	}

	if conn.URI == "" {
		return nil, errors.New("Missing URI")
	}

	if conn.Source != connector.VCS && conn.Source != connector.Local && conn.Source != connector.HTTP {
		return nil, errors.New("Source not supported")
	}

	if connectorForm.NameOverride != "" {
		conn.Name = connectorForm.NameOverride
	}

	switch conn.Source {
	case connector.VCS:
		err = downloadFromVCS(conn)
	case connector.HTTP:
		err = downloadFromHTTP(conn)
	case connector.Local:
		err = downloadFromLocal(conn)
	}

	if err != nil {
		return nil, err
	}
	// WIP Add cleanup if Addconnector fails
	return p.ConnectorRepository.AddConnector(conn)
}

func (p *ConnectorManager) DeleteConnector(name string) error {
	err := p.ConnectorRepository.DeleteConnector(name)
	if err != nil {
		log.Println(err)
		return errors.New("Unable to remove connector from database")
	}
	err = DeleteConnector(name)
	if err != nil {
		return errors.New("Unable to delete connector")
	}
	return nil
}

func (p *ConnectorManager) EditConnector(conn *connector.Connector) (*connector.Connector, error) {
	return p.ConnectorRepository.EditConnector(conn)
}

func (p *ConnectorManager) UpdateConnector(name string) error {
	conn, err := p.ConnectorRepository.GetConnector(name)
	if err != nil {
		return err
	}
	if conn.URI == "" {
		return errors.New("Missing URI")
	}

	if conn.Source != connector.VCS && conn.Source != connector.Local && conn.Source != connector.HTTP {
		return errors.New("Source not supported")
	}

	switch conn.Source {
	case connector.VCS:
		return updateFromVCS(conn)
	case connector.HTTP:
		return updateFromHTTP(conn)
	case connector.Local:
		return updateFromLocal(conn)
	}

	return nil
}

func (p *ConnectorManager) ListConnectors() ([]connector.Connector, int, error) {
	return p.ConnectorRepository.ListConnectors()
}

func (p *ConnectorManager) ListAccounts() ([]connector.Account, int, error) {
	return p.ConnectorRepository.ListAccounts()
}

func (p *ConnectorManager) GetAccount(id int32) (*connector.Account, error) {
	return p.ConnectorRepository.GetAccount(id)
}

func (p *ConnectorManager) AddAccount(acc *connector.Account) (*connector.Account, error) {
	return p.ConnectorRepository.AddAccount(acc)
}

func (p *ConnectorManager) EditAccount(account *connector.Account) (*connector.Account, error) {
	return p.ConnectorRepository.EditAccount(account)
}

func (p *ConnectorManager) DeleteAccount(id int32) error {
	return p.ConnectorRepository.DeleteAccount(id)
}

func (p *ConnectorManager) ListData(filters *types.Filter) ([]connector.Data, int, error) {
	return p.ConnectorRepository.ListData(filters)
}

func (p *ConnectorManager) GetData(id int32) (*connector.Data, error) {
	return p.ConnectorRepository.GetData(id)
}

func (p *ConnectorManager) EditData(data *connector.Data) (*connector.Data, error) {
	return p.ConnectorRepository.EditData(data)
}
