package types

type ConnectorForm struct {
	NameOverride string `form:"name_override"`
	Config       string `form:"config"`
}

type Filter struct {
	Search     string   `query:"search"`
	Page       int      `query:"page"`
	Accounts   []int    `query:"account_id"`
	Connectors []string `query:"connector"`
	OrderBy    string   `query:"order_by"`
	Sort       string   `query:"sort"`
}
