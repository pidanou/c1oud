package types

type ConnectorForm struct {
	NameOverride string `form:"name_override" json:"name_override"`
	Config       string `form:"config" json:"config"`
}

type DataFilter struct {
	Search     string   `query:"search"`
	Page       int      `query:"page"`
	Limit      int      `query:"limit"`
	Accounts   []int    `query:"account_id"`
	Connectors []string `query:"connector"`
	OrderBy    string   `query:"order_by"`
	Sort       string   `query:"sort"`
}

type SyncInfoFilter struct {
	Connectors []string `query:"connector"`
	Accounts   []int    `query:"account_id"`
	Success    *bool    `query:"success"`
	Page       int      `query:"page"`
	Limit      int      `query:"limit"`
	OrderBy    string   `query:"order_by"`
	Sort       string   `query:"sort"`
}
