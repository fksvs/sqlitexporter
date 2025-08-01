package sqlitexporter

import (
	"context"
	"database/sql"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const (
	CreateLogsTableSQL = `
		CREATE TABLE IF NOT EXISTS %s (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			ResourceSchemaURL TEXT,
			ScopeSchemaURL TEXT,
			ObservedTimestamp TEXT,
			Timestamp TEXT,
			SeverityNumber INTEGER,
			SeverityText TEXT,
			Body TEXT,
			TraceID TEXT,
			SpanID TEXT,
			EventName TEXT,
			ResourceAttrs BLOB,
			LogAttrs BLOB
		);
	`

	InsertLogsTableSQL = `
		INSERT INTO %s (
			ResourceSchemaURL,
			ScopeSchemaURL,
			ObservedTimestamp,
			Timestamp,
			SeverityNumber,
			SeverityText,
			Body,
			TraceID,
			SpanID,
			EventName,
			ResourceAttrs,
			LogAttrs
		) VALUES (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12);
	`
)

type logsExporter struct {
	client *sql.DB
	logger *zap.Logger
	cfg    *Config
}

func newLogsExporter(logger *zap.Logger, cfg *Config) (*logsExporter, error) {
	db, err := cfg.BuildDB()
	if err != nil {
		return nil, err
	}

	return &logsExporter{
		client: db,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (e *logsExporter) Start(ctx context.Context, _ component.Host) error {
	e.logger.Info("starting sqlitexporter",
		zap.String("filename", e.cfg.DatabaseFilename),
		zap.String("logs table", e.cfg.LogsTableName),
	)
	return CreateLogsTable(ctx, e.cfg, e.client)
}

func (e *logsExporter) Shutdown(_ context.Context) error {
	e.logger.Info("shutting down sqlitexporter",
		zap.String("filename", e.cfg.DatabaseFilename),
		zap.String("logs table", e.cfg.LogsTableName),
	)
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

func (e *logsExporter) pushLogData(ctx context.Context, record plog.Logs) error {
	return nil
}

func CreateLogsTable(ctx context.Context, cfg *Config, client *sql.DB) error {
	str := fmt.Sprintf(CreateLogsTableSQL, cfg.LogsTableName)

	_, err := client.ExecContext(ctx, str)
	if err != nil {
		return err
	}

	return nil
}
