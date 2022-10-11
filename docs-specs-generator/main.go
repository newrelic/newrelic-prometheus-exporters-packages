// Package synthesis implements the mapping of metrics into NR entities
// The entity synthesis mapping logic is based on this project (https://github.com/newrelic-experimental/entity-synthesis-definitions).
// The definition of rules are the same to the ones defined in the definition.yaml files of the mentioned repo.
//
// Copyright 2021 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	dto "github.com/prometheus/client_model/go"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	flagPromMetricsPath = "prom_metrics_path"
	flagSpecPath        = "spec_path"
	flagDocsPath        = "docs_path"
)

//go:embed input/template.tmpl
var docTemplateContent string

func main() {
	promMetrics := flag.String(flagPromMetricsPath, "input/metrics.prom", "Path to the prom metrics file")
	specFile := flag.String(flagSpecPath, "output/specFile.yaml", "Path to the output spec file")
	docsFile := flag.String(flagDocsPath, "output/docsFile.html", "Path to the output docs file")
	flag.Parse()

	c, err := loadConfig()
	if err != nil {
		log.Errorf("loading config: %v", err)
		return
	}

	metrics, err := getPromMetrics(*promMetrics)
	if err != nil {
		log.Errorf("getting prom metrics: %v", err)
		return
	}

	sp := generateSpecFile(c, metrics, *specFile)

	generateDocFile(sp, *docsFile)

	return
}

func loadConfig() (Config, error) {
	c := Config{}
	cfg := viper.New()
	cfg.AddConfigPath("./input")
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")

	err := cfg.ReadInConfig()
	if err != nil {
		return c, err
	}

	err = cfg.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func getPromMetrics(filename string) ([]Metric, error) {
	var metricsCap int
	mfs, err := readMetrics(filename)
	if err != nil {
		return nil, err
	}
	for _, mf := range mfs {
		_, ok := supportedMetricTypes[mf.GetType()]
		if !ok {
			continue
		}
		metricsCap += len(mf.Metric)
	}

	metrics := make([]Metric, 0, metricsCap)
	for mname, mf := range mfs {
		ntype := mf.GetType()
		mtype, ok := supportedMetricTypes[ntype]
		if !ok {
			continue
		}
		for _, m := range mf.GetMetric() {
			var value interface{}
			var nrType metricType
			switch ntype {
			case io_prometheus_client.MetricType_UNTYPED:
				value = m.GetUntyped().GetValue()
				nrType = metricType_GAUGE
			case io_prometheus_client.MetricType_COUNTER:
				value = m.GetCounter().GetValue()
				nrType = metricType_COUNTER
			case io_prometheus_client.MetricType_GAUGE:
				value = m.GetGauge().GetValue()
				nrType = metricType_GAUGE
			case io_prometheus_client.MetricType_SUMMARY:
				value = m.GetSummary()
				nrType = metricType_SUMMARY
			case io_prometheus_client.MetricType_HISTOGRAM:
				value = m.GetHistogram()
				nrType = metricType_HISTOGRAM
			default:
				log.Printf("\"metric type not supported: %s\"", mtype)
				continue
			}
			attrs := map[string]interface{}{}
			for _, l := range m.GetLabel() {
				attrs[l.GetName()] = l.GetValue()
			}
			metrics = append(
				metrics,
				Metric{
					name:        mname,
					metricType:  nrType,
					value:       value,
					attributes:  attrs,
					description: mf.GetHelp(),
				},
			)
		}
	}
	return sortMetrics(metrics), nil
}

// Get scrapes the given URL and decodes the retrieved payload.
func readMetrics(filename string) (MetricFamiliesByName, error) {
	mfs := MetricFamiliesByName{}
	r, err := os.Open(filename)
	defer r.Close()
	if err != nil {
		return nil, fmt.Errorf("reading metrics %w", err)
	}
	d := expfmt.NewDecoder(r, expfmt.FmtText)
	for {
		var mf dto.MetricFamily
		if err := d.Decode(&mf); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		mfs[mf.GetName()] = mf
	}
	return mfs, nil
}

func sortMetrics(metrics []Metric) []Metric {
	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].name < metrics[j].name
	})
	return metrics
}
