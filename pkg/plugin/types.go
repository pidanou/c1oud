package plugin

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
)

const (
	VCS   string = "vcs"
	Local string = "local"
	HTTP  string = "http"
)

type Plugin struct {
	Name           string `db:"name"`
	Source         string `db:"source"`
	URI            string `db:"uri"`
	InstallCommand string `db:"install_command"`
	UpdateCommand  string `db:"update_command"`
	Command        string `db:"command"`
}

type Account struct {
	ID      int32  `db:"id"`
	Plugin  string `db:"plugin"`
	Name    string `db:"name"`
	Options string `db:"options"`
}

type Data struct {
	ID           int32    `db:"id"`
	RemoteID     string   `db:"remote_id"`
	Plugin       string   `db:"plugin"`
	ResourceName string   `db:"resource_name"`
	URI          string   `db:"uri"`
	Metadata     Metadata `db:"metadata"`
}

type Metadata map[string]string

func (m Metadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Metadata, invalid type")
	}
	return json.Unmarshal(bytes, &m)
}

func (m Metadata) Value() (driver.Value, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println("Cannot marshal to json", err)
		return "{}", nil
	}
	return string(bytes), nil
}
