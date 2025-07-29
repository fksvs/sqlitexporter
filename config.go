package sqlitexporter

type Config struct {
	DatabaseFilename string `mapstruct:"database_filename"`
	LogsTableName    string `mapstructure:"logs_table_name"`
}

func (cfg *Config) Validate() error {
	return nil
}
