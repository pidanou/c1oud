package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/pidanou/c1-core/internal/connectormanager"
	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/internal/migrations"
	"github.com/pidanou/c1-core/internal/repositories"
	"github.com/pidanou/c1-core/internal/server"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

func main() {
	dirname, err := os.UserHomeDir()
	path := filepath.Join(dirname, ".c1")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}
	constants.Envs["C1_DIR"] = path

	var db = &sqlx.DB{}
	dbEngine := os.Getenv("C1_DB_ENGINE")

	connectorManager := &connectormanager.ConnectorManager{}
	if dbEngine != "sqlite" {
		db, err = setupPostgresDB()
		connectorManager = connectormanager.NewConnectorManager(repositories.NewPostgresRepository(db))
	} else {
		db, err = setupSQLiteDB()
		connectorManager = connectormanager.NewConnectorManager(repositories.NewSQLiteRepository(db))
	}
	if err != nil {
		log.Fatalf("Error booting connectormanager: %v", err)
	}
	defer db.Close()

	app := &server.Server{ConnectorManager: connectorManager, DB: db}
	app.Start(*connectorManager)
}

func setupPostgresDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", os.Getenv("C1_POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	var m *migrate.Migrate

	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB")
		return nil, err
	}

	useOS := os.Getenv("env") == "dev"
	if useOS {
		m, err = migrate.NewWithDatabaseInstance(
			"file://migrations/postgres/scripts",
			"pgx",
			driver)
		if err != nil {
			log.Fatal("Error generating migration: ", err)
		}
	} else {
		mig, err := iofs.New(migrations.Migrations, "scripts/postgres")
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
		log.Println(err)
		m.Down()
		log.Fatal("Cannot update DB")
		return nil, err
	}
	return db, nil
}

func setupSQLiteDB() (*sqlx.DB, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot open DB")
	}
	dbPath := userDir + "/.c1/c1.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath) // Create the empty DB file
		if err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
	}
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Cannot open DB: %s", err)
	}
	var m *migrate.Migrate

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB")
		return nil, err
	}

	useOS := os.Getenv("env") == "dev"
	if useOS {
		m, err = migrate.NewWithDatabaseInstance(
			"file://migrations/sqlite/scripts",
			"sqlite",
			driver)
		if err != nil {
			log.Fatal("Error generating migration: ", err)
		}
	} else {
		mig, err := iofs.New(migrations.Migrations, "scripts/sqlite")
		if err != nil {
			log.Fatal("Error getting migrations files: ", err)
		}

		m, err = migrate.NewWithInstance(
			"iofs", mig,
			"sqlite",
			driver)
		if err != nil {
			log.Fatal("Error generating migration: ", err)
		}
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println(err)
		m.Down()
		log.Fatal("Cannot update DB")
		return nil, err
	}
	return db, nil
}
