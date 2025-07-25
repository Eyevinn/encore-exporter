package main

import (
	"context"
	"errors"
	"github/Eyevinn/encore-exporter/cmd/exporter/telemetry"
	"github/Eyevinn/encore-exporter/internal/encore"
	"github/Eyevinn/encore-exporter/internal/logger"
	"log/slog"
	"os"
	"time"
)

type config struct {
	pollingIntervalMs int
}

func main() {
	config := &config{
		pollingIntervalMs: 2000, // poll API every 500 milliseconds
	}

	otelShutdown, err := telemetry.SetupOtelSdk(
		context.Background(),
		"encore-exporter",
		"1.0.0",
		true,
	)

	if err != nil {
		logger.Error("Failed to setup OpenTelemetry SDK", slog.String("error", err.Error()))
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	encoreMetrics, err := telemetry.NewEncoreMetrics()
	if err != nil {
		logger.Error("Failed to create Encore metrics", slog.String("error", err.Error()))
		os.Exit(1)
	}

	encoreHandler, err := encore.NewEncoreHandler(encoreMetrics)

	if err != nil {
		logger.Error("Error during initialization", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Creating ticker for polling API")
	ticker := time.NewTicker(
		time.Duration(config.pollingIntervalMs) * time.Millisecond,
	)

	defer ticker.Stop()

	for range ticker.C {
		logger.Debug("Polling API for metrics")
		// Here you would typically call your API to collect metrics
		// and update the EncoreMetrics instance accordingly.
		err := encoreHandler.GetEncoreMetrics(time.Duration(config.pollingIntervalMs) * time.Millisecond)
		if err != nil {
			logger.Error("Failed to get Encore metrics", slog.String("error", err.Error()))
			continue
		}
		logger.Debug("Metrics collected and updated")
	}
}
