package prometheus

import (
	dto "github.com/prometheus/client_model/go"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

// Metric represents a Prometheus metric.
// https://prometheus.io/docs/concepts/data_model/
type Metric struct {
	Name        string
	Value       metricValue
	MetricType  metricType
	Attributes  Set
	Description string
}

type metricValue interface{}
type metricType string
type MetricFamiliesByName map[string]dto.MetricFamily
type Set map[string]interface{}

var supportedMetricTypes = map[io_prometheus_client.MetricType]string{
	io_prometheus_client.MetricType_COUNTER:   "counter",
	io_prometheus_client.MetricType_GAUGE:     "gauge",
	io_prometheus_client.MetricType_HISTOGRAM: "histogram",
	io_prometheus_client.MetricType_SUMMARY:   "summary",
	io_prometheus_client.MetricType_UNTYPED:   "untyped",
}

//nolint:golint
const (
	MetricType_COUNTER   metricType = "count"
	MetricType_GAUGE     metricType = "gauge"
	MetricType_SUMMARY   metricType = "summary"
	MetricType_HISTOGRAM metricType = "histogram"
)
