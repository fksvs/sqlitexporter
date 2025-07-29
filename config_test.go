package sqlitexporter

import "testing"

func TestValidate(t *testing.T) {
	var cfg Config

	cfg.DatabaseFilename = "database.db"
	cfg.LogsTableName = "logs"

	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}
