# otel-grpc-healthcheck

Health Check extension enables a gRPC url that can be probed to check the status of the OpenTelemetry Collector. This extension was created when the OpenTelemetry collector did not support gRPC health check. If there is a better alternative please use it. This extension was created solely to be able to deploy the OpenTelemetry Collector in gateway-mode behind an AWS Application Load Balancer (ALB). The extension depends on the http health check extension.

To use the extension the OpenTelemetry Collector needs to be built using the OpenTelemetry Collector Builder.

## The OpenTelemetry Collector Builder Configuration

Example, adjust to your needs:
```yaml
dist:
  name: otel-collector-custom
  description: Custom collector
  output_path: ./bin
  otelcol_version: 0.130.0

exporters:
  - gomod: go.opentelemetry.io/collector/exporter/debugexporter v0.130.0

extensions:
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.130.0 # Required
  - gomod: github.com/jammymalina/otel-grpc-healthcheck v0.130.0
    import: github.com/jammymalina/otel-grpc-healthcheck
    name: grpc_health_check

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.130.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.130.0

receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.130.0

providers:
  - gomod: go.opentelemetry.io/collector/confmap/provider/envprovider v1.36.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/fileprovider v1.36.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpprovider v1.36.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.36.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.36.0
```

## The OpenTelemetry Collector Configuration

Example, adjust to your needs:
```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"
      http:
        endpoint: "0.0.0.0:4318"

processors:
  resource:
    attributes:
      - key: service.version
        value: ${env:APP_VERSION}
        action: upsert
      - key: deployment.environment
        value: ${env:ENV}
        action: upsert

  batch:
    send_batch_max_size: 1000
    send_batch_size: 100
    timeout: 10s

exporters:
  debug:
    verbosity: detailed
    sampling_initial: 5
    sampling_thereafter: 200

extensions:
  health_check:
    endpoint: "localhost:13133"
  grpc_health_check:
    grpc:
      endpoint: "0.0.0.0:13134"
    health_check_http_endpoint: "http://localhost:13133"

service:
  extensions: [health_check, grpc_health_check]
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [resource, batch]
      exporters: [debug]
    traces:
      receivers: [otlp]
      processors: [resource, batch]
      exporters: [debug]
```
