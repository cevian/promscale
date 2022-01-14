// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license.

package api

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/timescale/promscale/pkg/util"
)

var metrics *Metrics

type Metrics struct {
	// Using the first word in struct to ensure proper alignment in 32-bit systems.
	// Reference: https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	LastRequestUnixNano          int64
	LeaderGauge                  prometheus.Gauge
	ReceivedSamples              prometheus.Counter
	ReceivedMetadata             prometheus.Counter
	FailedSamples                prometheus.Counter
	FailedMetadata               prometheus.Counter
	IngestedSamples              prometheus.Counter
	SentMetadata                 prometheus.Counter
	SentBatchDuration            prometheus.Histogram
	ReceivedQueries              prometheus.Counter
	ExecutedQueries              prometheus.Counter
	TimedOutQueries              prometheus.Counter
	FailedQueries                prometheus.Counter
	QueryRemoteReadBatchDuration prometheus.Histogram
	ExemplarQueryDuration        prometheus.Histogram
	QueryDuration                prometheus.Histogram
	InvalidReadReqs              prometheus.Counter
	InvalidWriteReqs             prometheus.Counter
	InvalidQueryReqs             prometheus.Counter
	HTTPRequestDuration          *prometheus.HistogramVec
}

// InitMetrics sets up and returns the Prometheus metrics which Promscale exposes.
// This needs to be set before calling objects from the api package.
func InitMetrics() *Metrics {
	if metrics != nil {
		return metrics
	}
	metrics = createMetrics()
	prometheus.MustRegister(
		metrics.LeaderGauge,
		metrics.ReceivedSamples,
		metrics.ReceivedMetadata,
		metrics.ReceivedQueries,
		metrics.IngestedSamples,
		metrics.SentMetadata,
		metrics.FailedSamples,
		metrics.FailedMetadata,
		metrics.FailedQueries,
		metrics.ExecutedQueries,
		metrics.TimedOutQueries,
		metrics.InvalidReadReqs,
		metrics.InvalidWriteReqs,
		metrics.SentBatchDuration,
		metrics.QueryRemoteReadBatchDuration,
		metrics.QueryDuration,
		metrics.ExemplarQueryDuration,
		metrics.HTTPRequestDuration,
	)

	return metrics
}

func createMetrics() *Metrics {
	return &Metrics{
		LeaderGauge: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: util.PromNamespace,
				Name:      "current_leader",
				Help:      "Shows current election leader status",
			},
		),
		ReceivedSamples: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "received_samples_total",
				Help:      "Total number of received samples.",
			},
		),
		ReceivedMetadata: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "received_metadata_total",
				Help:      "Total number of received metadata.",
			},
		),
		FailedSamples: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "failed_samples_total",
				Help:      "Total number of processed samples which failed on send to remote storage.",
			},
		),
		FailedMetadata: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "failed_metadata_total",
				Help:      "Total number of processed metadata which failed on send to remote storage.",
			},
		),
		IngestedSamples: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "ingested_samples_total",
				Help:      "Total number of processed samples sent to remote storage.",
			},
		),
		SentMetadata: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "sent_metadata_total",
				Help:      "Total number of processed metadata sent to remote storage.",
			},
		),
		SentBatchDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: util.PromNamespace,
				Name:      "sent_batch_duration_seconds",
				Help:      "Duration of sample batch send calls to the remote storage.",
				Buckets:   prometheus.DefBuckets,
			},
		),
		LastRequestUnixNano: time.Now().UnixNano(),
		QueryRemoteReadBatchDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: util.PromNamespace,
				Subsystem: "metrics",
				Name:      "query_remote_read_batch_duration_seconds",
				Help:      "Duration of query batch read calls to the remote storage.",
				Buckets:   []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 50, 100, 250, 500, 1000, 2500},
			},
		),
		QueryDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: util.PromNamespace,
				Subsystem: "metrics",
				Name:      "query_duration_seconds",
				Help:      "Duration of query batch read calls to the PromQL engine.",
				Buckets:   []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 50, 100, 250, 500, 1000, 2500},
			},
		),
		ExemplarQueryDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: util.PromNamespace,
				Name:      "exemplar_query_duration_seconds",
				Help:      "Duration of exemplar query read calls to the database.",
				Buckets:   prometheus.DefBuckets,
			},
		),
		FailedQueries: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Subsystem: "metrics",
				Name:      "queries_failed_total",
				Help:      "Total number of queries which failed on send to remote storage.",
			},
		),
		InvalidReadReqs: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "invalid_read_requests",
				Help:      "Total number of remote read requests with invalid metadata.",
			},
		),
		InvalidWriteReqs: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "invalid_write_requests",
				Help:      "Total number of remote write requests with invalid metadata.",
			},
		),
		InvalidQueryReqs: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "invalid_query_requests",
				Help:      "Total number of invalid query requests with invalid metadata.",
			},
		),
		ReceivedQueries: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "received_queries_total",
				Help:      "Total number of received queries.",
			},
		),
		ExecutedQueries: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "executed_queries_total",
				Help:      "Total number of successfully executed queries.",
			},
		),
		TimedOutQueries: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: util.PromNamespace,
				Name:      "queries_timed_out_total",
				Help:      "Total number of timed out queries.",
			},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: util.PromNamespace,
				Name:      "http_request_duration_ms",
				Help:      "Duration of HTTP request in milliseconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"path"},
		),
	}
}
