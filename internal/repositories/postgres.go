package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/internal/types"
	"github.com/pidanou/c1-core/pkg/connector"
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(DB *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{DB: DB}
}

func (p *PostgresRepository) ListConnectors() ([]connector.Connector, int, error) {
	var connectors []connector.Connector
	var count = 0
	query := `SELECT * FROM connectors`
	err := p.DB.Select(&connectors, query)
	if err != nil {
		return nil, count, err
	}
	query = `SELECT count(*) from connectors`
	err = p.DB.Get(&count, query)
	if err != nil {
		return nil, count, err
	}
	return connectors, count, nil
}

func (p *PostgresRepository) GetConnector(name string) (*connector.Connector, error) {
	var connector connector.Connector
	query := `SELECT * FROM connectors WHERE name = $1 LIMIT 1`
	err := p.DB.Get(&connector, query, name)
	if err != nil {
		return nil, err
	}
	return &connector, nil
}

func (p *PostgresRepository) AddConnector(plug *connector.Connector) (*connector.Connector, error) {
	query := `INSERT INTO connectors (name, description, source, uri, install_command, update_command, command) VALUES (:name, :description, :source, :uri, :install_command, :update_command, :command)`
	_, err := p.DB.NamedExec(query, plug)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

func (p *PostgresRepository) DeleteConnector(name string) error {
	query := `DELETE FROM connectors WHERE name = $1`
	_, err := p.DB.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) EditConnector(plug *connector.Connector) (*connector.Connector, error) {
	query := `UPDATE connectors set description = :description, source = :source, uri = :uri, install_command = :install_command, update_command = :update_command, command = :command WHERE name = :name`
	_, err := p.DB.NamedExec(query, plug)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

func (p *PostgresRepository) ListAccounts() ([]connector.Account, int, error) {
	var accounts []connector.Account
	var count = 0
	query := `SELECT * FROM accounts`
	err := p.DB.Select(&accounts, query)
	if err != nil {
		return nil, count, err
	}
	query = `SELECT count(*) from accounts`
	err = p.DB.Get(&count, query)
	if err != nil {
		return nil, count, err
	}
	return accounts, count, nil
}

func (p *PostgresRepository) GetAccount(id int32) (*connector.Account, error) {
	var account connector.Account
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1`
	err := p.DB.Get(&account, query, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (p *PostgresRepository) AddAccount(acc *connector.Account) (*connector.Account, error) {
	query := `INSERT INTO accounts (name, connector, options) VALUES (:name, :connector, :options)`
	_, err := p.DB.NamedExec(query, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (p *PostgresRepository) EditAccount(acc *connector.Account) (*connector.Account, error) {
	query := `UPDATE accounts set name = :name, connector = :connector, options = :options WHERE id = :id`
	_, err := p.DB.NamedExec(query, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (p *PostgresRepository) DeleteAccount(id int32) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := p.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) AddData(data []connector.Data) error {
	query := `INSERT INTO data (account_id, remote_id, resource_name, connector, uri, metadata) VALUES (:account_id, :remote_id, :resource_name, :connector, :uri, :metadata) ON CONFLICT (remote_id, account_id) 
	DO UPDATE SET 
		resource_name = EXCLUDED.resource_name,
		connector = EXCLUDED.connector,
		uri = EXCLUDED.uri,
		metadata = EXCLUDED.metadata`
	_, err := p.DB.NamedExec(query, data)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) ListData(filters *types.Filter) ([]connector.Data, int, error) {
	var data []connector.Data
	var count = 0
	query := `SELECT * FROM data WHERE 1=1`
	countQuery := `SELECT count(*) FROM data WHERE 1=1`
	query, countQuery, args, err := p.buildQuery(query, countQuery, filters)
	if err != nil {
		return nil, count, err
	}
	err = p.DB.Select(&data, query, args...)
	if err != nil {
		return nil, count, err
	}
	err = p.DB.Get(&count, countQuery, args...)
	if err != nil {
		return nil, count, err
	}
	return data, count, nil
}

func (p *PostgresRepository) GetData(id int32) (*connector.Data, error) {
	data := &connector.Data{}
	query := `SELECT * FROM data WHERE id = $1`
	err := p.DB.Get(data, query, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PostgresRepository) EditData(data *connector.Data) (*connector.Data, error) {
	query := `UPDATE data SET notes = :notes WHERE id = :id`
	_, err := p.DB.NamedExec(query, data)
	if err != nil {
		return nil, err
	}
	return p.GetData(data.ID)
}

func (p *PostgresRepository) buildQuery(baseQuery string, countQuery string, filters *types.Filter) (string, string, []interface{}, error) {
	var args []interface{}
	if filters == nil {
		baseQuery += fmt.Sprintf(" ORDER BY account_id ASC, resource_name ASC LIMIT %v", constants.PageSize)
		return baseQuery, countQuery, []interface{}{}, nil
	}
	if filters.Search != "" {
		baseQuery += fmt.Sprint(" AND (resource_name ILIKE ? OR metadata ILIKE ? OR notes ILIKE ?)")
		countQuery += fmt.Sprint(" AND (resource_name ILIKE ? OR metadata ILIKE ? OR notes ILIKE ?)")
		args = append(args, "%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%", "%")
	}
	if filters.Accounts != nil {
		queryPart, argsPart, _ := sqlx.In(" AND account_id in (?)", filters.Accounts)
		baseQuery += queryPart
		countQuery += queryPart
		args = append(args, argsPart...)
	}
	if filters.Connectors != nil {
		queryPart, argsPart, _ := sqlx.In(" AND connector in (?)", filters.Connectors)
		baseQuery += queryPart
		countQuery += queryPart
		args = append(args, argsPart...)
	}
	if filters.OrderBy != "" && isValidOrderBy(filters.OrderBy) {
		baseQuery += fmt.Sprintf("%v ORDER BY %s", baseQuery, filters.OrderBy)
	} else {
		baseQuery += " ORDER BY account_id ASC, resource_name ASC"
	}
	if filters.Sort != "" && (filters.Sort == "ASC" || filters.Sort == "DESC") {
		baseQuery += fmt.Sprint(" %v", filters.Sort)
	}
	baseQuery += " LIMIT 50"
	if filters.Page != 0 {
		baseQuery += fmt.Sprintf(" Offset %v", (filters.Page-1)*constants.PageSize)
	}
	baseQuery = p.DB.Rebind(baseQuery)
	countQuery = p.DB.Rebind(countQuery)
	return baseQuery, countQuery, args, nil
}

func isValidOrderBy(orderBy string) bool {
	allowedColumns := map[string]bool{
		"account_id":    true,
		"resource_name": true,
	}
	return allowedColumns[orderBy]
}
