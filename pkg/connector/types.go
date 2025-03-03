package connector

import "fmt"

const (
	VCS   string = "VCS"
	Local string = "Local"
	HTTP  string = "HTTP"
)

type Connector struct {
	Name           string `db:"name" form:"name" json:"name"`
	Description    string `db:"description" form:"description" json:"description"`
	Source         string `db:"source" form:"source" json:"source"`
	URI            string `db:"uri" form:"uri" json:"uri"`
	InstallCommand string `db:"install_command" form:"install_command" json:"install_command"`
	UpdateCommand  string `db:"update_command" form:"update_command" json:"update_command"`
	Command        string `db:"command" form:"command" json:"command"`
}

func (c Connector) ToJSON() string {
	return fmt.Sprintf(`{"name": "%s", "description": "%s", "source": "%s", "uri": "%s", "install_command": "%s", "update_command": "%s", "command": "%s"}`, c.Name, c.Description, c.Source, c.URI, c.InstallCommand, c.UpdateCommand, c.Command)
}

type Account struct {
	ID        int32  `db:"id"`
	Connector string `db:"connector" form:"connector"`
	Name      string `db:"name" form:"name"`
	Options   string `db:"options" form:"options"`
}

type Data struct {
	ID           int32  `db:"id"`
	AccountID    int32  `db:"account_id"`
	RemoteID     string `db:"remote_id"`
	Connector    string `db:"connector"`
	ResourceName string `db:"resource_name"`
	URI          string `db:"uri"`
	Metadata     string `db:"metadata"`
	Notes        string `db:"notes" form:"notes"`
}

type SyncInfo struct {
	ID        int32  `db:"id"`
	Connector string `db:"connector"`
	AccountID string `db:"account_id"`
	Metadata  string `db:"metadata"`
	CreatedAt string `db:"created_at"`
	Success   bool   `db:"success"`
}
