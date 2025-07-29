package sqlitexporter

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
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

// TODO: write table initializer
func (e *logsExporter) Start(ctx context.Context, _ component.Host) error {
	e.logger.Info("starting sqlitexporter",
		zap.String("filename", e.cfg.DatabaseFilename),
		zap.String("logs table", e.cfg.LogsTableName),
	)
	return nil
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
