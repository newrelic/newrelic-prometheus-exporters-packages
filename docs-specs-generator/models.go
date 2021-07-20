package main

import (
	dto "github.com/prometheus/client_model/go"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

//nolint:golint
const (
	metricType_COUNTER   metricType = "count"
	metricType_GAUGE     metricType = "gauge"
	metricType_SUMMARY   metricType = "summary"
	metricType_HISTOGRAM metricType = "histogram"
)

var supportedMetricTypes = map[io_prometheus_client.MetricType]string{
	io_prometheus_client.MetricType_COUNTER:   "counter",
	io_prometheus_client.MetricType_GAUGE:     "gauge",
	io_prometheus_client.MetricType_HISTOGRAM: "histogram",
	io_prometheus_client.MetricType_SUMMARY:   "summary",
	io_prometheus_client.MetricType_UNTYPED:   "untyped",
}

type Specs struct {
	SpecVersion                  string    `yaml:"specVersion"`
	OwningTeam                   string    `yaml:"owningTeam"`
	IntegrationName              string    `yaml:"integrationName"`
	HumanReadableIntegrationName string    `yaml:"humanReadableIntegrationName"`
	Entities                     []*Entity `yaml:"entities"`
}

// Metrics
type MetricSpec struct {
	Name              string `yaml:"name"`
	Type              string `yaml:"type"`
	DefaultResolution int    `yaml:"defaultResolution"`
	Unit              string `yaml:"unit"`
	Description       string
	Dimensions        []Dimension `yaml:"dimensions"`
}

// InternalAttributes
type InternalAttribute struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// Entity
type Entity struct {
	EntityType         string              `yaml:"entityType"`
	Metrics            []*MetricSpec       `yaml:"metrics"`
	InternalAttributes []InternalAttribute `yaml:"internalAttributes"`
	IgnoredAttributes  []string            `yaml:"ignoredAttributes"`
}

type Dimension struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// Metric represents a Prometheus metric.
// https://prometheus.io/docs/concepts/data_model/
type Metric struct {
	name        string
	value       metricValue
	metricType  metricType
	attributes  Set
	description string
}

// Config
type Config struct {
	IntegrationName              string       `mapstructure:"integrationName"`
	HumanReadableIntegrationName string       `mapstructure:"humanReadableIntegrationName"`
	Definitions                  []Definition `mapstructure:"definitions"` // List of conditions used find metrics that match.
}

// EntityRule contains rules to synthesis an entity
type Definition struct {
	EntityType string      `mapstructure:"type"`
	Conditions []Condition `mapstructure:"conditions"` // List of rules used to determining if a metric belongs to this entity.
}

// Condition decides whether a metric is suitable to be included into an entity based on the metric attributes and metric name.
type Condition struct {
	Attribute string `mapstructure:"attribute"`
	Prefix    string `mapstructure:"prefix"`
	Value     string `mapstructure:"value"`
}

type metricValue interface{}
type metricType string
type MetricFamiliesByName map[string]dto.MetricFamily
type Set map[string]interface{}
