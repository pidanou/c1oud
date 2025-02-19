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

func (s *Server) Start() {
	err := s.setupDB()
	if err != nil {
		log.Panic(err)
	}

	e := echo.New()
	e.HideBanner = true

	useOS := os.Getenv("env") == "dev"
	assetHandler := http.FileServer(getFileSystem(useOS, ui.StaticFiles))

	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.GET("/", GetHome)
	e.Logger.Fatal(e.Start(":1323"))
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
