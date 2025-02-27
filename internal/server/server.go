package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pidanou/c1-core/internal/connectormanager"
	"github.com/pidanou/c1-core/internal/ui"
)

type Server struct {
	DB               *sqlx.DB
	ConnectorManager *connectormanager.ConnectorManager
}

func getFileSystem(useOS bool, embededFiles embed.FS) http.FileSystem {
	if useOS {
		return http.FS(os.DirFS("internal/ui/static"))
	}

	fsys, err := fs.Sub(embededFiles, "internal/ui/static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func (s *Server) Start(connManager connectormanager.ConnectorManager) error {
	h := &Handler{ConnectorManager: connManager}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())

	useOS := os.Getenv("ENV") == "dev"
	assetHandler := http.FileServer(getFileSystem(useOS, ui.StaticFiles))

	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/data")
	})
	e.GET("/data", h.GetDataPage)
	e.GET("/accounts", h.GetAccountsPage)
	e.GET("/account/new", h.GetNewAccountPage)
	e.POST("/account", h.PostAccount)
	e.GET("/connectors", h.GetConnectorsPage)
	e.GET("/connector/new", h.GetNewConnectorPage)
	e.POST("/connector", h.PostConnector)

	partials := e.Group("/partials")

	partials.GET("/data/:id/edit", h.GetEditDataRow)
	partials.GET("/data", h.GetData)
	partials.PUT("/data/:id", h.PutData)
	partials.POST("/data/sync", h.PostDataSync)

	partials.GET("/connector/:name/edit", h.GetEditConnectorRow)
	partials.GET("/connector/:name", h.GetConnectorRow)
	partials.DELETE("/connector/:name", h.DeleteConnector)
	partials.PUT("/connector/:name", h.PutConnector)
	partials.POST("/connector/:name/update", h.PostConnectorUpdate)

	partials.GET("/account/:id/edit", h.GetEditAccountRow)
	partials.GET("/account/:id", h.GetAccountRow)
	partials.DELETE("/account/:id", h.DeleteAccount)
	partials.PUT("/account/:id", h.PutAccount)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":7777"
	}
	e.Logger.Fatal(e.Start(port))
	return nil
}
