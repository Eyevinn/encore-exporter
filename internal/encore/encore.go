package encore

import (
	"encoding/json"
	"github/Eyevinn/encore-exporter/cmd/exporter/telemetry"
	"github/Eyevinn/encore-exporter/internal/logger"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type EncoreJob struct {
	Id                  string    `json:"id,omitempty"`
	ExternalId          string    `json:"externalId,omitempty"`
	Profile             string    `json:"profile"`
	OutputFolder        string    `json:"outputFolder"`
	BaseName            string    `json:"baseName"`
	Status              string    `json:"status,omitempty"`
	ProgressCallbackUri string    `json:"progressCallbackUri,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
}

type EncoreHandler interface {
	GetEncoreMetrics(interval time.Duration) error
}

type EncoreClient struct {
	client        *http.Client
	encoreMetrics telemetry.EncoreMetricsCollector
	encoreUrl     string
}

// TODO: Implement
func NewEncoreHandler(encoreMetrics telemetry.EncoreMetricsCollector) (EncoreHandler, error) {
	// Implementation of the Encore client creation
	// This is a placeholder; actual implementation will depend on the Encore API
	return &EncoreClient{
		client:        &http.Client{},
		encoreMetrics: encoreMetrics,
	}, nil
}

func (ec *EncoreClient) getEncoreJobs(tickInterval time.Duration) ([]EncoreJob, error) {
	jobRequest, err := http.NewRequest("GET", ec.encoreUrl+"/encoreJobs", nil)
	if err != nil {
		return nil, err
	}
	qps := jobRequest.URL.Query()
	qps.Add("size", "100") // high number to get all jobs that could've been created during the polling interval
	qps.Add("sort", "createdAt,desc")
	jobRequest.URL.RawQuery = qps.Encode()
	resp, err := ec.client.Do(jobRequest)
	if err != nil {
		logger.Error("Failed to get Encore jobs", slog.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to get Encore jobs", slog.String("status", resp.Status))
		return nil, err
	}
	var encoreResp encoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&encoreResp); err != nil {
		logger.Error("Failed to decode Encore jobs response", slog.String("error", err.Error()))
		return nil, err
	}
	return encoreResp.Jobs(), nil
}

func (ec *EncoreClient) GetEncoreMetrics(interval time.Duration) error {
	jobs, err := ec.getEncoreJobs(interval)
	if err != nil {
		return err
	}
	jobs = filterOldJobs(jobs, interval)
	if len(jobs) == 0 {
		logger.Debug("No Encore jobs found in the last interval", slog.Duration("interval", interval))
		return nil
	}
	ec.IncrementMetrics(jobs)
	return nil
}

func filterOldJobs(jobs []EncoreJob, interval time.Duration) []EncoreJob {
	new := slices.Collect(func(yield func(EncoreJob) bool) {
		for _, job := range jobs {
			if job.CreatedAt.After(time.Now().Add(-interval)) {
				yield(job)
			}
		}
	})
	return new
}

// TODO: Verify this AI slop works
func (ec *EncoreClient) IncrementMetrics(jobs []EncoreJob) {
	for _, job := range jobs {
		switch job.Status {
		case "QUEUED":
			ec.encoreMetrics.EncoreJobQueue(1)
		case "FAILED":
			ec.encoreMetrics.EncoreJobFail(1)
		case "SUCCESS":
			ec.encoreMetrics.EncoreJobSuccess(1)
		case "IN_PROGRESS":
			ec.encoreMetrics.EncoreJobInProgress(1)
		default:
			logger.Warn("Unknown job status", slog.String("jobId", job.Id), slog.String("status", job.Status))
		}
	}
}

type encoreResponse struct {
	Embedded embedded `json:"_embedded"`
}

type embedded struct {
	EncoreJobs []EncoreJob `json:"encoreJobs"`
}

func (e *encoreResponse) Jobs() []EncoreJob {
	if e == nil {
		return nil
	}
	return e.Embedded.EncoreJobs
}
