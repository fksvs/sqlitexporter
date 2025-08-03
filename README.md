# SQLite Exporter for OpenTelemetry Collector

A custom exporter component for the OpenTelemetry Collector that exports log data to SQLite databases.

## Features

- **Lightweight Storage**: Uses SQLite for efficient, file-based log storage
- **OpenTelemetry Integration**: Seamlessly integrates with the OpenTelemetry Collector ecosystem
- **Configurable**: Supports custom database filenames and table names

## Installation and Usage

### Prerequisites

- Go 1.21 or later
- OpenTelemetry Collector Builder (`ocb`)

### Installing OpenTelemetry Collector Builder

Install the OpenTelemetry Collector Builder tool:

```bash
go install go.opentelemetry.io/collector/cmd/builder@latest
```

### Building the Collector

1. Create a builder configuration file (`builder-config.yaml`):

```yaml
dist:
  name: otelcol-custom
  description: Custom OpenTelemetry Collector distribution with SQLite exporter
  output_path: ./otelcol-custom
  version: 1.0.0

exporters:
  - gomod: github.com/fksvs/sqlitexporter v1.0.0

receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.130.1

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.130.1
  - gomod: go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.130.1

extensions:
  - gomod: go.opentelemetry.io/collector/extension/healthcheckextension v0.130.1
  - gomod: go.opentelemetry.io/collector/extension/pprofextension v0.130.1
```

2. Build the custom collector:

```bash
builder --config builder-config.yaml
```

3. Run the built collector:

```bash
./otelcol-custom/otelcol-custom --config collector-config.yaml
```

### Example Test Builder Configuration

For testing purposes, you can use this minimal configuration (`test-builder-config.yaml`):

```yaml
dist:
  name: otelcol-test
  description: Test build with SQLite exporter
  output_path: ./otelcol-test
  version: 0.1.0

exporters:
  - gomod: github.com/fksvs/sqlitexporter main

receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.130.1

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.130.1

extensions:
  - gomod: go.opentelemetry.io/collector/extension/healthcheckextension v0.130.1
```

Build the test collector:

```bash
builder --config test-builder-config.yaml
```

## Configuration

### Exporter Configuration

Add the SQLite exporter to your OpenTelemetry Collector configuration:

```yaml

exporters:
  sqlite:
    database_filename: "telemetry.db"
    logs_table_name: "logs"

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [sqlite]
```

### Configuration Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `database_filename` | string | Yes | - | Path to the SQLite database file |
| `logs_table_name` | string | Yes | - | Name of the table to store log data |

### Complete Example Configuration

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  
  memory_limiter:
    limit_mib: 512

exporters:
  sqlite:
    database_filename: "/var/log/otel/telemetry.db"
    logs_table_name: "application_logs"

extensions:
  health_check:
    endpoint: 0.0.0.0:13133
  pprof:
    endpoint: 0.0.0.0:1777

service:
  extensions: [health_check, pprof]
  pipelines:
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [sqlite]
```

## Contributing

Pull requests are welcome. For bug fixes and small improvements, please submit a pull request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is free software; you can redistribute it and/or modify it under the terms of the GPLv3 license. See LICENSE for details.