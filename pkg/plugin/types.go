package plugin

const (
	VCS   string = "VCS"
	Local string = "Local"
	HTTP  string = "HTTP"
)

type Plugin struct {
	Name           string `db:"name" form:"name" json:"name"`
	Source         string `db:"source" form:"source" json:"source"`
	URI            string `db:"uri" form:"uri" json:"uri"`
	InstallCommand string `db:"install_command" form:"install_command" json:"install_command"`
	UpdateCommand  string `db:"update_command" form:"update_command" json:"update_command"`
	Command        string `db:"command" form:"command" json:"command"`
}

type Account struct {
	ID      int32  `db:"id"`
	Plugin  string `db:"plugin" form:"plugin"`
	Name    string `db:"name" form:"name"`
	Options string `db:"options" form:"options"`
}

type Data struct {
	ID           int32  `db:"id"`
	AccountID    int32  `db:"account_id"`
	RemoteID     string `db:"remote_id"`
	Plugin       string `db:"plugin"`
	ResourceName string `db:"resource_name"`
	URI          string `db:"uri"`
	Metadata     string `db:"metadata"`
	Tags         string `db:"tags" form:"tags"`
	Notes        string `db:"notes" form:"notes"`
}
