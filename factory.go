package grpc_health_check

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/extension"
)

var (
	Type               = component.MustNewType("grpc_health_check")
	ExtensionStability = component.StabilityLevelBeta
)

const (
	defaultEndpoint                = "0.0.0.0:13134"
	defaultHealthcheckHttpEndpoint = "http://localhost:13133"
	defaultStartPeriod             = 30 * time.Second
	defaultInterval                = 5 * time.Second
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
		Grpc: configgrpc.ServerConfig{NetAddr: confignet.AddrConfig{
			Endpoint:  defaultEndpoint,
			Transport: "tcp",
		}},
		HealthCheckHttpEndpoint: defaultHealthcheckHttpEndpoint,
		StartPeriod:             defaultStartPeriod,
		Interval:                defaultInterval,
	}
	return cfg
}

func createExtension(_ context.Context, params extension.Settings, cfg component.Config) (extension.Extension, error) {
	config := cfg.(*Config)

	return newServer(*config, params.TelemetrySettings), nil
}
