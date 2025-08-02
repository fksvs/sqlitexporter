package sqlitexporter

import (
	"context"
	"fmt"
	"testing"
)

func TestExporterLogs(t *testing.T) {}

func TestCreateLogsTable(t *testing.T) {
	cfg := Config{
		DatabaseFilename: "./testdata/test_database.db",
		LogsTableName:    "logs",
	}

	db, err := cfg.BuildDB()
	if err != nil {
		t.Errorf("failed to build database, database filename: %s, logs table name: %s, err: %v",
			cfg.DatabaseFilename, cfg.LogsTableName, err)
	}

	err = CreateLogsTable(context.Background(), &cfg, db)
	if err != nil {
		t.Errorf("failed to create logs table, database filename: %s, logs table name: %s, err: %v",
			cfg.DatabaseFilename, cfg.LogsTableName, err)
	}

	_, err = db.ExecContext(context.Background(), fmt.Sprintf("SELECT ID FROM %s;", cfg.LogsTableName))
	if err != nil {
		t.Errorf("failed to select ID column database filename: %s, logs table name: %s, err: %v",
			cfg.DatabaseFilename, cfg.LogsTableName, err)
	}
}
