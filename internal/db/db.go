package db

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const driverName = "postgres"

type Settings struct {
	Host          string
	Port          string
	User          string
	Password      string
	DbName        string
	Schema        string
	AppName       string
	MaxIdleConns  int
	MaxOpenConns  int
	MigrateSchema bool
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Connect(settings Settings) (*sqlx.DB, error) {
	addr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s application_name=%s timezone=UTC",
		settings.Host,
		settings.Port,
		settings.User,
		settings.Password,
		settings.DbName,
		settings.Schema,
		settings.AppName,
	)

	sqlxDB, err := sqlx.Connect(driverName, addr)
	if err != nil {
		return nil, err
	}
	sqlxDB.SetMaxIdleConns(settings.MaxIdleConns)
	sqlxDB.SetMaxOpenConns(settings.MaxOpenConns)

	if err := goose.SetDialect(driverName); err != nil {
		return nil, err
	}

	if settings.MigrateSchema {
		goose.SetBaseFS(embedMigrations)
		if err := goose.Up(sqlxDB.DB, "migrations"); err != nil {
			return nil, err
		}
	}

	return sqlxDB, nil
}
