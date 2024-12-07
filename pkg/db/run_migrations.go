package db

import (
	"database/sql"

	"github.com/codeboris/music-lib/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbApi *sql.DB, cfg config.ConfigMigrate) error {
	migrator, err := migrate.New(cfg.MigrationPath, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
