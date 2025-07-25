package telemetry

import "errors"

type EncoreMetricsCollector interface {
	EncoreJobQueue(increment int64)
	EncoreJobFail(increment int64)
	EncoreJobSuccess(increment int64)
	EncoreJobInProgress(increment int64)
}

type EncoreMetrics struct {
	encoreJobFails      MetricsCollector
	encoreJobSuccess    MetricsCollector
	encoreJobInProgress MetricsCollector
	encoreJobQueue      MetricsCollector
}

func NewEncoreMetrics() (*EncoreMetrics, error) {
	em := EncoreMetrics{}

	encoreJobFail, encoreJobFailErr := NewCounter(
		"encore_job_fail",
		"Number of Encore jobs that failed",
	)
	encoreJobQueue, encoreJobQueueErr := NewCounter(
		"encore_job_queue",
		"Number of Encore jobs in queue",
	)
	encoreJobSuccess, encoreJobSuccessErr := NewCounter(
		"encore_job_success",
		"Number of Encore jobs that succeeded",
	)
	encoreJobInProgress, encoreJobInProgressErr := NewCounter(
		"encore_job_in_progress",
		"Number of Encore jobs in progress",
	)
	em.encoreJobFails = encoreJobFail
	em.encoreJobQueue = encoreJobQueue
	em.encoreJobSuccess = encoreJobSuccess
	em.encoreJobInProgress = encoreJobInProgress

	return &em, errors.Join(
		encoreJobFailErr,
		encoreJobQueueErr,
		encoreJobSuccessErr,
		encoreJobInProgressErr,
	)
}

func (e *EncoreMetrics) EncoreJobQueue(increment int64) {
	e.encoreJobQueue.IncrementCounter(increment)
}

func (e *EncoreMetrics) EncoreJobFail(increment int64) {
	e.encoreJobFails.IncrementCounter(increment)
}

func (e *EncoreMetrics) EncoreJobSuccess(increment int64) {
	e.encoreJobSuccess.IncrementCounter(increment)
}
func (e *EncoreMetrics) EncoreJobInProgress(increment int64) {
	e.encoreJobInProgress.IncrementCounter(increment)
}
