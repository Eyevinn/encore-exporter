package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceName("eyevinn/encore-exporter"),
			semconv.ServiceVersion("1.0.0"),
		))
}

func SetupOtelSdk(
	ctx context.Context,
	serviceName string,
	serviceVersion string,
	httpExport bool,
) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, f := range shutdownFuncs {
			err = errors.Join(err, f(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	var metricExporter metric.Exporter
	var meErr error

	if httpExport {
		metricExporter, meErr = otlpmetrichttp.New(ctx)
	} else {
		metricExporter, meErr = stdoutmetric.New()
	}
	if meErr != nil {
		handleErr(meErr)
		return
	}
	resource, err := newResource()

	if err != nil {
		handleErr(err)
		return
	}

	// Set up meter provider
	meterProvider := newMeterProvider(resource, metricExporter)
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	// Create a new propagator for context propagation
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newMeterProvider(res *resource.Resource, me metric.Exporter) *metric.MeterProvider {
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(me, metric.WithInterval(10*time.Second)),
		),
	)
	return meterProvider
}
