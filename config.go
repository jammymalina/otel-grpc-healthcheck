package grpc_health_check

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
)

type Config struct {
	Grpc                    configgrpc.GRPCServerSettings `mapstructure:"grpc"`
	HealthCheckHttpEndpoint string                        `mapstructure:"health_check_http_endpoint"`
}

var _ component.Config = (*Config)(nil)
