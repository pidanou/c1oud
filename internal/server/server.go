package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/pidanou/c1-core/internal/migrations"
	"github.com/pidanou/c1-core/internal/pluginmanager"
	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/internal/ui"
)

type Server struct {
	DB            *sqlx.DB
	PluginManager *pluginmanager.PluginManager
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

func (s *Server) Start() error {
	err := s.setupDB()
	if err != nil {
		log.Panic(err)
		return err
	}

	h := &Handler{PluginManager: *pluginmanager.NewPluginManager(repositories.NewPostgresRepository(s.DB))}

	e := echo.New()
	e.HideBanner = true

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
	e.GET("/plugins", h.GetPluginsPage)
	e.GET("/plugin/new", h.GetNewPluginPage)
	e.POST("/plugin", h.PostPlugin)

	partials := e.Group("/partials")

	partials.GET("/data/:id/edit", h.GetEditDataRow)
	partials.GET("/data", h.GetData)
	partials.PUT("/data/:id", h.PutData)
	partials.POST("/data/sync", h.PostDataSync)

	partials.GET("/plugin/:name/edit", h.GetEditPluginRow)
	partials.GET("/plugin/:name", h.GetPluginRow)
	partials.DELETE("/plugin/:name", h.DeletePlugin)
	partials.PUT("/plugin/:name", h.PutPlugin)
	partials.POST("/plugin/:name/update", h.PostPluginUpdate)

	partials.GET("/account/:id/edit", h.GetEditAccountRow)
	partials.GET("/account/:id", h.GetAccountRow)
	partials.DELETE("/account/:id", h.DeleteAccount)
	partials.PUT("/account/:id", h.PutAccount)

	e.Logger.Fatal(e.Start(":1323"))
	return nil
}

func (s *Server) setupDB() error {
	var m *migrate.Migrate

	driver, err := pgx.WithInstance(s.DB.DB, &pgx.Config{})
	if err != nil {
		return err
	}

	useOS := os.Getenv("env") == "dev"
	if useOS {
		m, err = migrate.NewWithDatabaseInstance(
			"file://migrations/scripts",
			"pgx",
			driver)
		if err != nil {
			log.Fatal("Error generating migration: ", err)
		}
	} else {
		mig, err := iofs.New(migrations.Migrations, "scripts")
		if err != nil {
			log.Fatal("Error getting migrations files: ", err)
		}

		m, err = migrate.NewWithInstance(
			"iofs", mig,
			"pgx",
			driver)
		if err != nil {
			log.Fatal("Error generating migration: ", err)
		}
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println(err)
		m.Down()
		return err
	}
	return nil
}
