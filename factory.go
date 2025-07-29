package sqlitexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

func NewFactory() exporter.Factory               {}
func createDefaultConfig() component.Config      {}
func createLogsExporter() (exporter.Logs, error) {}
