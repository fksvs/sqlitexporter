package sqlitexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/xexporter"
)

func NewFactory() exporter.Factory {
	return xexporter.NewFactory(
		component.MustNewType("sqlite"),
		createDefaultConfig,
		xexporter.WithLogs(createLogsExporter, component.StabilityLevelDevelopment))
}

func createDefaultConfig() component.Config {
	return &Config{
		DatabaseFilename: "database",
		LogsTableName:    "logs",
	}
}

func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	config := cfg.(*Config)

	exporter, err := newLogsExporter(set.Logger, config)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		exporter.pushLogData,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}
