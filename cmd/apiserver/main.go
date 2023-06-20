package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rasulov-emirlan/topenergy-interview/config"
	"github.com/rasulov-emirlan/topenergy-interview/internal/domains"
	"github.com/rasulov-emirlan/topenergy-interview/internal/storage/redis"
	"github.com/rasulov-emirlan/topenergy-interview/internal/transport/httprest"
	"github.com/rasulov-emirlan/topenergy-interview/pkg/health"
	"github.com/rasulov-emirlan/topenergy-interview/pkg/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logging.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	repo, err := redis.NewRepoCombiner(ctx, cfg)
	if err != nil {
		log.Fatal("failed to initialize redis repo", logging.Error("err", err))
	}

	doms, err := domains.NewDomainCombiner(domains.CommonDependencies{Log: log}, domains.TasksDependencies{Repo: repo.Tasks()})
	if err != nil {
		log.Fatal("failed to initialize domains", logging.Error("err", err))
	}

	jeagerShutdown, err := JaegerTraceProvider(cfg)
	if err != nil {
		log.Fatal("failed to initialize jaeger trace provider", logging.Error("err", err))
	}

	srv := httprest.NewServer(cfg)
	go func() {
		if err := srv.Start(log, doms, []health.Checker{repo.Check}); err != nil {
			log.Fatal("failed to start http server", logging.Error("err", err))
		}
	}()

	log.Info("server started", logging.String("port", cfg.Server.Port))

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info("shutting down server")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("failed to stop http server", logging.Error("err", err))
	}

	if err := repo.Close(); err != nil {
		log.Fatal("failed to close redis repo", logging.Error("err", err))
	}

	if err := jeagerShutdown(ctx); err != nil {
		log.Fatal("failed to shutdown jaeger trace provider", logging.Error("err", err))
	}

	log.Info("server stopped")
}

func JaegerTraceProvider(cfg config.Config) (func(context.Context) error, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JeagerURL)))
	if err != nil {
		return nil, err
	}

	envKey := "production"
	if cfg.Flags.DevMode {
		envKey = "development"
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(httprest.ServiceName),
			semconv.DeploymentEnvironmentKey.String(envKey),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp.Shutdown, nil
}
