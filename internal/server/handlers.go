package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/internal/pluginmanager"
	"github.com/pidanou/c1-core/internal/types"
	"github.com/pidanou/c1-core/internal/ui"
	"github.com/pidanou/c1-core/pkg/plugin"
)

type Handler struct {
	PluginManager pluginmanager.PluginManager
}

func Render(ctx echo.Context, statusCode int, t ...templ.Component) error {
	fmt.Println()
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	component := ""

	for _, tpl := range t {
		if err := tpl.Render(ctx.Request().Context(), buf); err != nil {
			return err
		}
		component = buf.String()
	}

	return ctx.HTML(statusCode, component)
}

func (h *Handler) GetDataPage(c echo.Context) error {
	filter := &types.Filter{}
	c.Bind(filter)
	page := filter.Page
	if page == 0 {
		page = 1
	}
	if filter.Accounts == nil {
		filter.Accounts = []int{}
	}
	if filter.Plugins == nil {
		filter.Plugins = []string{}
	}
	// TODO: tags manager
	// if filter.Tags == nil {
	//   filter.Tags = []string{}
	// }
	plugins, _, err := h.PluginManager.ListPlugins()
	if err != nil {
		plugins = []plugin.Plugin{}
	}

	accounts, _, err := h.PluginManager.ListAccounts()
	if err != nil {
		log.Println(err)
		accounts = []plugin.Account{}
	}

	data, count, err := h.PluginManager.ListData(filter)
	if err != nil {
		log.Println(err)
		data = []plugin.Data{}
	}

	return Render(c, http.StatusOK, ui.DataPage(accounts, data, plugins, page, (count+constants.PageSize-1)/constants.PageSize))
}

func (h *Handler) GetData(c echo.Context) error {
	filter := &types.Filter{}
	c.Bind(filter)
	page := filter.Page
	if page == 0 {
		page = 1
	}
	accounts, _, err := h.PluginManager.ListAccounts()
	if err != nil {
		log.Println(err)
	}

	data, count, err := h.PluginManager.ListData(filter)
	if err != nil {
		log.Println(err)
	}

	return Render(c, http.StatusOK, ui.DataTableBody(accounts, data), ui.OOB(ui.DataPagination(len(data) == constants.PageSize, page > 1, page, (count+constants.PageSize-1)/constants.PageSize), "outerHTML:.pagination"))
}

func (h *Handler) GetEditDataRow(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	data, err := h.PluginManager.GetData(int32(idInt))
	if err != nil {
		log.Println(err)
	}

	acc, err := h.PluginManager.GetAccount(data.AccountID)
	if err != nil {
		log.Println(err)
	}
	if acc == nil {
		acc = &plugin.Account{}
	}

	return Render(c, http.StatusOK, ui.DataRow(data, acc, true))
}

func (h *Handler) PutData(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var data = &plugin.Data{}
	err = c.Bind(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	data.ID = int32(idInt)
	data, err = h.PluginManager.EditData(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	acc, err := h.PluginManager.GetAccount(data.AccountID)
	if err != nil {
		log.Println(err)
	}
	if acc == nil {
		acc = &plugin.Account{}
	}
	return Render(c, http.StatusOK, ui.DataRow(data, acc, false))
}

func (h *Handler) PostDataSync(c echo.Context) error {
	accounts, _, err := h.PluginManager.ListAccounts()
	if err != nil {
		log.Println(err)
	}

	accountIDs := []int32{}
	for _, acc := range accounts {
		accountIDs = append(accountIDs, acc.ID)
	}

	err = h.PluginManager.Execute(accountIDs)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	data, _, err := h.PluginManager.ListData(nil)
	if err != nil {
		log.Println(err)
	}

	return Render(c, http.StatusOK, ui.DataTableBody(accounts, data))
}

func (h *Handler) GetAccountsPage(c echo.Context) error {
	accounts, _, err := h.PluginManager.ListAccounts()
	if err != nil {
		log.Println(err)
		return Render(c, http.StatusOK, ui.AccountsPage(nil))
	}
	return Render(c, http.StatusOK, ui.AccountsPage(accounts))
}

func (h *Handler) GetPluginsPage(c echo.Context) error {
	plugins, _, err := h.PluginManager.ListPlugins()
	if err != nil {
		log.Println(err)
		return Render(c, http.StatusInternalServerError, ui.PluginsPage(nil))
	}
	return Render(c, http.StatusOK, ui.PluginsPage(plugins))
}

func (h *Handler) PostPlugin(c echo.Context) error {
	pluginForm := &types.PluginForm{}
	err := c.Bind(pluginForm)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to read config", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to read config")
	}

	if pluginForm.Config == "" {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to read config", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to read config")
	}

	_, err = h.PluginManager.InstallPlugin(pluginForm)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to install plugin", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to install plugin")
	}

	c.Response().Header().Set("Hx-Redirect", `/plugins`)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetNewPluginPage(c echo.Context) error {
	return Render(c, http.StatusOK, ui.NewPluginPage())
}

func (h *Handler) GetEditPluginRow(c echo.Context) error {
	name := c.Param("name")
	plug, err := h.PluginManager.GetPlugin(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := true
	return Render(c, http.StatusOK, ui.PluginRow(plug, edit))
}

func (h *Handler) GetPluginRow(c echo.Context) error {
	name := c.Param("name")
	plug, err := h.PluginManager.GetPlugin(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := false
	return Render(c, http.StatusOK, ui.PluginRow(plug, edit))
}

func (h *Handler) PutPlugin(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.String(http.StatusInternalServerError, "No plugin name")
	}
	var plug = &plugin.Plugin{}
	err := c.Bind(plug)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	plug.Name = name
	plug, err = h.PluginManager.EditPlugin(plug)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return Render(c, http.StatusOK, ui.PluginRow(plug, false))
}

func (h *Handler) PostPluginUpdate(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "No plugin selected", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "No plugin selected")
	}
	err := h.PluginManager.UpdatePlugin(name)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	plug, err := h.PluginManager.GetPlugin(name)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return Render(c, http.StatusOK, ui.PluginRow(plug, false))
}

func (h *Handler) DeletePlugin(c echo.Context) error {
	name := c.Param("name")
	err := h.PluginManager.DeletePlugin(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(200, "")
}

func (h *Handler) GetEditAccountRow(c echo.Context) error {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	account, err := h.PluginManager.GetAccount(int32(intId))
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	plugins, _, err := h.PluginManager.ListPlugins()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := true
	return Render(c, http.StatusOK, ui.AccountRow(account, plugins, edit))
}

func (h *Handler) GetAccountRow(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	account, err := h.PluginManager.GetAccount(int32(idInt))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := false
	return Render(c, http.StatusOK, ui.AccountRow(account, nil, edit))
}

func (h *Handler) DeleteAccount(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = h.PluginManager.DeleteAccount(int32(idInt))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(200, "")
}

func (h *Handler) GetNewAccountPage(c echo.Context) error {
	plugins, _, err := h.PluginManager.ListPlugins()
	if err != nil {
		log.Println(err)
	}
	return Render(c, http.StatusOK, ui.NewAccountPage(plugins))
}

func (h *Handler) PostAccount(c echo.Context) error {
	account := &plugin.Account{}
	err := c.Bind(account)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable add account", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to add account")
	}

	if account.Plugin == "" || account.Name == "" {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Missing field", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Missing field")
	}

	_, err = h.PluginManager.AddAccount(account)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to add account", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to add account")
	}

	c.Response().Header().Set("Hx-Redirect", `/accounts`)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) PutAccount(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var acc = &plugin.Account{}
	err = c.Bind(acc)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	acc.ID = int32(idInt)
	acc, err = h.PluginManager.EditAccount(acc)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return Render(c, http.StatusOK, ui.AccountRow(acc, nil, false))
}
