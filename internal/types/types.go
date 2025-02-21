package types

type PluginForm struct {
	NameOverride string `form:"name_override"`
	Config       string `form:"config"`
}
