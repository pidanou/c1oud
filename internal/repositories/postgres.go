package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/pidanou/c1-core/pkg/plugin"
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(DB *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{DB: DB}
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

func (p *PostgresRepository) GetAccount(id int) (*plugin.Account, error) {
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

func (p *PostgresRepository) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE name = $1`
	_, err := p.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) AddData(data []plugin.Data) error {
	query := `INSERT INTO data (remote_id, resource_name, plugin, uri, metadata) VALUES (:remote_id, :resource_name, :plugin, :uri, :metadata)`
	_, err := p.DB.NamedExec(query, data)
	if err != nil {
		return err
	}
	return nil
}
