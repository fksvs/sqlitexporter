package sqlitexporter

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DatabaseFilename string `mapstructure:"database_filename"`
	LogsTableName    string `mapstructure:"logs_table_name"`
}

func (cfg *Config) Validate() error {
	if cfg.DatabaseFilename == "" || cfg.LogsTableName == "" {
		return fmt.Errorf("database filename and Logs table name must be provided")
	}

	return nil
}

func (cfg *Config) BuildDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DatabaseFilename)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
