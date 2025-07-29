package sqlitexporter

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type LogsExporter struct {
	client *sql.DB
	logger *zap.Logger
	cfg    *Config
}

func newLogsExporter(logger *zap.Logger, cfg *Config) (*LogsExporter, error)      {}
func (e *LogsExporter) Start(ctx context.Context, _ component.Host) error         {}
func (e *LogsExporter) Shutdown(_ context.Context) error                          {}
func (e *LogsExporter) pushLogRecord(ctx context.Context, record plog.Logs) error {}
