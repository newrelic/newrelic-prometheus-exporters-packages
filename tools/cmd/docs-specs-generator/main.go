package main

import (
	_ "embed"
	"flag"

	"github.com/newrelic/newrelic-prometheus-exporters-packages/tools/src/docs"
	"github.com/newrelic/newrelic-prometheus-exporters-packages/tools/src/prometheus"
	"github.com/newrelic/newrelic-prometheus-exporters-packages/tools/src/specs"

	log "github.com/sirupsen/logrus"
)

const (
	flagPromMetricsPath = "prom_metrics_path"
	flagSpecPath        = "spec_path"
	flagDocsPath        = "docs_path"
)

func main() {
	promMetrics := flag.String(flagPromMetricsPath, "powerdns.prom", "Path to the prom metrics file")
	specFile := flag.String(flagSpecPath, "output/specFile.yaml", "Path to the output spec file")
	docsFile := flag.String(flagDocsPath, "output/docsFile.html", "Path to the output docs file")
	flag.Parse()

	c, err := specs.LoadConfig()
	if err != nil {
		log.Errorf("loading config: %v", err)
		return
	}

	metrics, err := prometheus.GetPromMetrics(*promMetrics)
	if err != nil {
		log.Errorf("getting prom metrics: %v", err)
		return
	}

	sp := specs.GenerateSpecFile(c, metrics, *specFile)

	docs.GenerateDocFile(sp, *docsFile)

	return
}
