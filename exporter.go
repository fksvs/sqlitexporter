package sqlitexporter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
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
			ResourceAttrs TEXT,
			LogAttrs TEXT
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
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
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
	sql := fmt.Sprintf(InsertLogsTableSQL, e.cfg.LogsTableName)
	stmt, err := e.client.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < record.ResourceLogs().Len(); i++ {
		resourceLogs := record.ResourceLogs().At(i)
		resourceSchemaURL := resourceLogs.SchemaUrl()
		resourceAttrs, err := json.Marshal(attributesToMap(resourceLogs.Resource().Attributes()))
		if err != nil {
			return nil
		}

		for j := 0; j < resourceLogs.ScopeLogs().Len(); j++ {
			scopeLogs := resourceLogs.ScopeLogs().At(j)
			scopeSchemaURL := scopeLogs.SchemaUrl()

			for k := 0; k < scopeLogs.LogRecords().Len(); k++ {
				logRecord := scopeLogs.LogRecords().At(k)
				observedTimestamp := logRecord.ObservedTimestamp().String()
				timestamp := logRecord.Timestamp().String()
				severityNumber := logRecord.SeverityNumber()
				severityText := logRecord.SeverityText()
				body := logRecord.Body().AsString()
				traceID := logRecord.TraceID().String()
				spanID := logRecord.SpanID().String()
				eventName := logRecord.EventName()
				logAttrs, err := json.Marshal(attributesToMap(logRecord.Attributes()))
				if err != nil {
					return nil
				}

				_, err = stmt.ExecContext(ctx,
					resourceSchemaURL,
					scopeSchemaURL,
					observedTimestamp,
					timestamp,
					severityNumber,
					severityText,
					body,
					traceID,
					spanID,
					eventName,
					resourceAttrs,
					logAttrs,
				)
				if err != nil {
					return err
				}
			}
		}
	}

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

func attributesToMap(attributes pcommon.Map) map[string]interface{} {
	result := make(map[string]interface{}, attributes.Len())
	attributes.Range(func(k string, v pcommon.Value) bool {
		switch v.Type() {
		case pcommon.ValueTypeStr:
			result[k] = v.Str()
		case pcommon.ValueTypeInt:
			result[k] = v.Int()
		case pcommon.ValueTypeDouble:
			result[k] = v.Double()
		case pcommon.ValueTypeBool:
			result[k] = v.Bool()
		case pcommon.ValueTypeBytes:
			result[k] = v.Bytes().AsRaw()
		default:
			result[k] = v.AsString()
		}
		return true
	})
	return result
}
