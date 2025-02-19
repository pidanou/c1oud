package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/internal/pluginmanager"
	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/internal/server"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dirname, err := os.UserHomeDir()
	path := filepath.Join(dirname, ".c1")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}
	constants.Envs["C1_DIR"] = path

	db, err := sqlx.Open("pgx", os.Getenv("C1_POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repositories.NewPostgresRepository(db)
	pluginManager := pluginmanager.NewPluginManager(repo)

	app := &server.Server{PluginManager: pluginManager, DB: db}
	app.Start()
}
