package grpc_health_check

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/extension"
)

const (
	Type               = "grpc_health_check"
	ExtensionStability = component.StabilityLevelBeta
)

const (
	defaultEndpoint               = "0.0.0.0:13134"
	defaultHealtcheckHttpEndpoint = "localhost:13134"
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		Type,
		createDefaultConfig,
		createExtension,
		ExtensionStability,
	)
}

func createDefaultConfig() component.Config {
	cfg := &Config{
		Grpc:                    configgrpc.GRPCServerSettings{},
		HealthCheckHttpEndpoint: defaultHealtcheckHttpEndpoint,
	}
	cfg.Grpc.NetAddr.Endpoint = defaultEndpoint
	return cfg
}

func createExtension(_ context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	config := cfg.(*Config)

	return newServer(*config, set.TelemetrySettings), nil
}
