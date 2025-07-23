package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alpineworks.io/ootel"
	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/config"
	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/logging"
	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/mqtt"
	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/timescale"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}

	slogLevel, err := logging.LogLevelToSlogLevel(logLevel)
	if err != nil {
		log.Fatalf("could not convert log level: %s", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})))
	c, err := config.NewConfig()
	if err != nil {
		slog.Error("could not create config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()

	exporterType := ootel.ExporterTypePrometheus
	if c.Local {
		exporterType = ootel.ExporterTypeOTLPGRPC
	}

	ootelClient := ootel.NewOotelClient(
		ootel.WithMetricConfig(
			ootel.NewMetricConfig(
				c.MetricsEnabled,
				exporterType,
				c.MetricsPort,
			),
		),
		ootel.WithTraceConfig(
			ootel.NewTraceConfig(
				c.TracingEnabled,
				c.TracingSampleRate,
				c.TracingService,
				c.TracingVersion,
			),
		),
	)

	shutdown, err := ootelClient.Init(ctx)
	if err != nil {
		slog.Error("could not create ootel client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(5 * time.Second))
	if err != nil {
		slog.Error("could not create runtime metrics", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = host.Start()
	if err != nil {
		slog.Error("could not create host metrics", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	mqttClient := mqtt.NewMQTTClient(c.MQTTAddress,
		mqtt.WithMQTTClientID(c.MQTTClientID),
		mqtt.WithMQTTKeepAlive(c.MQTTKeepAlive),
		mqtt.WithMQTTPingTimeout(c.MQTTPingTimeout),
		mqtt.WithMQTTConnectTimeout(c.MQTTConnectTimeout),
		mqtt.WithMQTTCleanSession(c.MQTTCleanSession),
	)
	err = mqttClient.Connect()
	if err != nil {
		slog.Error("could not connect to mqtt broker", slog.String("error", err.Error()))
		os.Exit(1)
	}

	timescaleClient, err := timescale.NewTimescaleClient(ctx, c.TimescaleConnString)
	if err != nil {
		slog.Error("could not create timescale client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer timescaleClient.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	processor := mqtt.NewProcessor(mqttClient, timescaleClient, ctx, c.Timezone)
	err = processor.Start(c.MQTTTopic, c.MQTTQoS)
	if err != nil {
		slog.Error("could not start processor", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("birdnet-mqtt-timescaledb started", slog.String("topic", c.MQTTTopic))

	<-sigChan
	slog.Info("shutdown signal received, disconnecting...")
	mqttClient.Disconnect()
}
