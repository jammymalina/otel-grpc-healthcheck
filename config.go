package grpc_health_check

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
)

type Config struct {
	Grpc                    configgrpc.ServerConfig `mapstructure:"grpc"`
	HealthCheckHttpEndpoint string                  `mapstructure:"health_check_http_endpoint"`
	StartPeriod             time.Duration           `mapstructure:"start_period"`
	Interval                time.Duration           `mapstructure:"interval"`
}

var _ component.Config = (*Config)(nil)
