package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pidanou/c1-core/pkg/plugin"
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(DB *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{DB: DB}
}

func (p *PostgresRepository) ListPlugins() ([]plugin.Plugin, error) {
	var plugins []plugin.Plugin
	query := `SELECT * FROM plugins`
	err := p.DB.Select(&plugins, query)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func (p *PostgresRepository) GetPlugin(name string) (*plugin.Plugin, error) {
	var plugin plugin.Plugin
	query := `SELECT * FROM plugins WHERE name = $1 LIMIT 1`
	err := p.DB.Get(&plugin, query, name)
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

func (p *PostgresRepository) AddPlugin(plug *plugin.Plugin) (*plugin.Plugin, error) {
	query := `INSERT INTO plugins (name, source, uri, install_command, update_command, command) VALUES (:name, :source, :uri, :install_command, :update_command, :command)`
	_, err := p.DB.NamedExec(query, plug)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

func (p *PostgresRepository) DeletePlugin(name string) error {
	query := `DELETE FROM plugins WHERE name = $1`
	_, err := p.DB.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) EditPlugin(plug *plugin.Plugin) (*plugin.Plugin, error) {
	query := `UPDATE plugins set source = :source, uri = :uri, install_command = :install_command, update_command = :update_command, command = :command WHERE name = :name`
	_, err := p.DB.NamedExec(query, plug)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

func (p *PostgresRepository) ListAccounts() ([]plugin.Account, error) {
	var accounts []plugin.Account
	query := `SELECT * FROM accounts`
	err := p.DB.Select(&accounts, query)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (p *PostgresRepository) GetAccount(id int32) (*plugin.Account, error) {
	var account plugin.Account
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1`
	err := p.DB.Get(&account, query, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (p *PostgresRepository) AddAccount(acc *plugin.Account) (*plugin.Account, error) {
	query := `INSERT INTO accounts (name, plugin) VALUES (:name, :plugin)`
	_, err := p.DB.NamedExec(query, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (p *PostgresRepository) EditAccount(acc *plugin.Account) (*plugin.Account, error) {
	query := `UPDATE accounts set name = :name, plugin = :plugin, options = :options WHERE id = :id`
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

func (p *PostgresRepository) AddData(data []plugin.Data) error {
	query := `INSERT INTO data (account_id, remote_id, resource_name, plugin, uri, metadata) VALUES (:account_id, :remote_id, :resource_name, :plugin, :uri, :metadata)`
	_, err := p.DB.NamedExec(query, data)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) ListData(limit, offset int, filters map[string]string) ([]plugin.Data, error) {
	var data []plugin.Data
	query := `SELECT * FROM data WHERE 1=1`
	query, args, err := p.buildQuery(query, filters)
	if err != nil {
		return nil, err
	}
	query += " ORDER BY id DESC LIMIT $" + fmt.Sprint(len(args)+1) + " OFFSET $" + fmt.Sprint(len(args)+2)
	args = append(args, limit, offset)
	err = p.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PostgresRepository) GetData(id int32) (*plugin.Data, error) {
	data := &plugin.Data{}
	query := `SELECT * FROM data WHERE id = $1`
	err := p.DB.Get(data, query, id)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	return data, nil
}

func (p *PostgresRepository) EditData(data *plugin.Data) (*plugin.Data, error) {
	query := `UPDATE data SET notes = :notes, tags = :tags WHERE id = :id`
	_, err := p.DB.NamedExec(query, data)
	if err != nil {
		return nil, err
	}
	return p.GetData(data.ID)
}

func (p *PostgresRepository) buildQuery(baseQuery string, filters map[string]string) (string, []interface{}, error) {
	var args []interface{}
	index := 1
	for key, value := range filters {
		baseQuery += fmt.Sprintf(" AND "+key+" = $%v", index)
		args = append(args, value)
		index++
	}
	return baseQuery, args, nil
}
