package grpc_health_check

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	client = http.Client{
		Timeout: 5 * time.Second,
	}
)

type grpcHealthCheckExtension struct {
	config   Config
	logger   *zap.Logger
	server   *grpc.Server
	stopCh   chan struct{}
	settings component.TelemetrySettings
}

func (gc *grpcHealthCheckExtension) Start(_ context.Context, host component.Host) error {
	gc.logger.Info("Starting grpc_health_check extension", zap.Any("config", gc.config))
	ln, err := gc.config.Grpc.ToListener()
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", gc.config.Grpc.NetAddr.Endpoint, err)
	}

	gc.server, err = gc.config.Grpc.ToServer(host, gc.settings)
	if err != nil {
		return err
	}

	gc.stopCh = make(chan struct{})
	hs := health.NewServer()

	// Register the health server with the gRPC server
	healthpb.RegisterHealthServer(gc.server, hs)
	reflection.Register(gc.server)

	go func() {
		for {
			status := healthpb.HealthCheckResponse_SERVING
			response, err := client.Get(gc.config.HealthCheckHttpEndpoint)
			if err != nil {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			} else if response.StatusCode < 200 || response.StatusCode >= 300 {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			}
			hs.SetServingStatus("", status)

			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		defer close(gc.stopCh)

		// The listener ownership goes to the server.
		if err = gc.server.Serve(ln); !errors.Is(err, http.ErrServerClosed) && err != nil {
			host.ReportFatalError(err)
		}
	}()

	return nil
}

func (gc *grpcHealthCheckExtension) Shutdown(context.Context) error {
	if gc.server == nil {
		return nil
	}
	gc.server.GracefulStop()
	if gc.stopCh != nil {
		<-gc.stopCh
	}
	return nil
}

func newServer(config Config, settings component.TelemetrySettings) *grpcHealthCheckExtension {
	return &grpcHealthCheckExtension{
		config:   config,
		logger:   settings.Logger,
		settings: settings,
	}
}
