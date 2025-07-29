package sqlitexporter

import "testing"

func TestNewFactory(t *testing.T) {}
func TestCreateDefaultConfigLogs(t *testing.T) {
	cfg := createDefaultConfig()

	if cfg == nil {
		t.Errorf("failed to create default config")
	}
}
func TestCreateLogsExporter(t *testing.T) {}
