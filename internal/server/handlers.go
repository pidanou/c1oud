package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pidanou/c1-core/internal/connectormanager"
	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/internal/types"
	"github.com/pidanou/c1-core/internal/ui"
	"github.com/pidanou/c1-core/pkg/connector"
)

type Handler struct {
	ConnectorManager connectormanager.ConnectorManager
}

func Render(ctx echo.Context, statusCode int, t ...templ.Component) error {
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
	if filter.Connectors == nil {
		filter.Connectors = []string{}
	}
	connectors, _, err := h.ConnectorManager.ListConnectors()
	if err != nil {
		connectors = []connector.Connector{}
	}

	accounts, _, err := h.ConnectorManager.ListAccounts()
	if err != nil {
		log.Println(err)
		accounts = []connector.Account{}
	}

	data, count, err := h.ConnectorManager.ListData(filter)
	if err != nil {
		log.Println(err)
		data = []connector.Data{}
	}

	return Render(c, http.StatusOK, ui.DataPage(accounts, data, connectors, page, (count+constants.PageSize-1)/constants.PageSize))
}

func (h *Handler) GetData(c echo.Context) error {
	filter := &types.Filter{}
	c.Bind(filter)
	page := filter.Page
	if page == 0 {
		page = 1
	}
	accounts, _, err := h.ConnectorManager.ListAccounts()
	if err != nil {
		log.Println(err)
	}

	data, count, err := h.ConnectorManager.ListData(filter)
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
	data, err := h.ConnectorManager.GetData(int32(idInt))
	if err != nil {
		log.Println(err)
	}

	acc, err := h.ConnectorManager.GetAccount(data.AccountID)
	if err != nil {
		log.Println(err)
	}
	if acc == nil {
		acc = &connector.Account{}
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

	var data = &connector.Data{}
	err = c.Bind(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	data.ID = int32(idInt)
	data, err = h.ConnectorManager.EditData(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	acc, err := h.ConnectorManager.GetAccount(data.AccountID)
	if err != nil {
		log.Println(err)
	}
	if acc == nil {
		acc = &connector.Account{}
	}
	return Render(c, http.StatusOK, ui.DataRow(data, acc, false))
}

func (h *Handler) PostDataSync(c echo.Context) error {
	form, _ := c.FormParams()

	accounts := form["account_id"]
	accountIDs := []int32{}
	for _, acc := range accounts {
		accInt, _ := strconv.Atoi(acc)
		accountIDs = append(accountIDs, int32(accInt))
	}

	err := h.ConnectorManager.Execute(accountIDs)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Failed to synced data", "type": "warning"}}`)
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Succesfully synced data", "type": "info"}}`)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetAccountsPage(c echo.Context) error {
	accounts, _, err := h.ConnectorManager.ListAccounts()
	if err != nil {
		log.Println(err)
		return Render(c, http.StatusOK, ui.AccountsPage(nil))
	}
	return Render(c, http.StatusOK, ui.AccountsPage(accounts))
}

func (h *Handler) GetConnectorsPage(c echo.Context) error {
	connectors, _, err := h.ConnectorManager.ListConnectors()
	if err != nil {
		log.Println(err)
		return Render(c, http.StatusInternalServerError, ui.ConnectorsPage(nil))
	}
	return Render(c, http.StatusOK, ui.ConnectorsPage(connectors))
}

func (h *Handler) PostConnector(c echo.Context) error {
	connectorForm := &types.ConnectorForm{}
	err := c.Bind(connectorForm)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to read config", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to read config")
	}

	if connectorForm.Config == "" {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to read config", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to read config")
	}

	_, err = h.ConnectorManager.InstallConnector(connectorForm)
	if err != nil {
		log.Println(err)
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable to install connector", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to install connector")
	}

	c.Response().Header().Set("Hx-Redirect", `/connectors`)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetNewConnectorPage(c echo.Context) error {
	return Render(c, http.StatusOK, ui.NewConnectorPage())
}

func (h *Handler) GetEditConnectorRow(c echo.Context) error {
	name := c.Param("name")
	plug, err := h.ConnectorManager.GetConnector(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := true
	return Render(c, http.StatusOK, ui.ConnectorRow(plug, edit))
}

func (h *Handler) GetConnectorRow(c echo.Context) error {
	name := c.Param("name")
	plug, err := h.ConnectorManager.GetConnector(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := false
	return Render(c, http.StatusOK, ui.ConnectorRow(plug, edit))
}

func (h *Handler) PutConnector(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.String(http.StatusInternalServerError, "No connector name")
	}
	var plug = &connector.Connector{}
	err := c.Bind(plug)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	plug.Name = name
	plug, err = h.ConnectorManager.EditConnector(plug)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return Render(c, http.StatusOK, ui.ConnectorRow(plug, false))
}

func (h *Handler) PostConnectorUpdate(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "No connector selected", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "No connector selected")
	}
	err := h.ConnectorManager.UpdateConnector(name)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	plug, err := h.ConnectorManager.GetConnector(name)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return Render(c, http.StatusOK, ui.ConnectorRow(plug, false))
}

func (h *Handler) DeleteConnector(c echo.Context) error {
	name := c.Param("name")
	err := h.ConnectorManager.DeleteConnector(name)
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
	account, err := h.ConnectorManager.GetAccount(int32(intId))
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	connectors, _, err := h.ConnectorManager.ListConnectors()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	edit := true
	return Render(c, http.StatusOK, ui.AccountRow(account, connectors, edit))
}

func (h *Handler) GetAccountRow(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", fmt.Sprintf(`{"add-toast": {"message": "%s", "type": "warning"}}`, err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	account, err := h.ConnectorManager.GetAccount(int32(idInt))
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
	err = h.ConnectorManager.DeleteAccount(int32(idInt))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(200, "")
}

func (h *Handler) GetNewAccountPage(c echo.Context) error {
	connectors, _, err := h.ConnectorManager.ListConnectors()
	if err != nil {
		log.Println(err)
	}
	return Render(c, http.StatusOK, ui.NewAccountPage(connectors))
}

func (h *Handler) PostAccount(c echo.Context) error {
	account := &connector.Account{}
	err := c.Bind(account)
	if err != nil {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Unable add account", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Unable to add account")
	}

	if account.Connector == "" || account.Name == "" {
		c.Response().Header().Set("Hx-Trigger", `{"add-toast": {"message": "Missing field", "type": "warning"}}`)
		return c.String(http.StatusInternalServerError, "Missing field")
	}

	_, err = h.ConnectorManager.AddAccount(account)
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

	var acc = &connector.Account{}
	err = c.Bind(acc)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	acc.ID = int32(idInt)
	acc, err = h.ConnectorManager.EditAccount(acc)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return Render(c, http.StatusOK, ui.AccountRow(acc, nil, false))
}
