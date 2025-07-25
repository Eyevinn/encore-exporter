package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter metric.Meter

type MetricsCollector interface {
	IncrementCounter(increment int64)
}

type Counter struct {
	calls metric.Int64Counter
}

func (c *Counter) IncrementCounter(increment int64) {
	c.calls.Add(context.Background(), increment)
}

func NewCounter(name string, description string) (*Counter, error) {
	if meter == nil {
		meter = otel.Meter("counters")
	}
	calls, counterErr := meter.Int64Counter(
		name,
		metric.WithDescription(description),
	)
	return &Counter{
		calls: calls,
	}, errors.Join(counterErr)
}
