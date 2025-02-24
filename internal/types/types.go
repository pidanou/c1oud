package types

type PluginForm struct {
	NameOverride string `form:"name_override"`
	Config       string `form:"config"`
}

type Filter struct {
	Limit    int      `query:"limit"`
	Offset   int      `query:"offset"`
	Accounts []int    `query:"account_id"`
	Plugins  []string `query:"plugin"`
	Tags     []string `query:"tag"`
	OrderBy  string   `query:"order_by"`
	Sort     string   `query:"sort"`
}
