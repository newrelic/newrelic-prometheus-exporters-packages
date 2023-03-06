package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/newrelic/newrelic-prometheus-exporters-packages/tools/src/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	counterSuffix_bucket = "_bucket"
	counterSuffix_count  = "_count"
	counterSuffix_total  = "_total"

	// Types required by newrelic https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration/#override-mapping
	nrCounter   = "counter"
	nrGauge     = "gauge"
	genericRule = `
- source_labels: [__name__]
  regex: ^%s$
  target_label: newrelic_metric_type
  replacement: "%s"
  action: replace`
)

var (
	promInput     string
	outputMetrics string
	outputRules   string
)

func main() {
	flag.StringVar(&promInput, "promInput", "./powerdns.prom", "source for the prometheus metrics")
	flag.StringVar(&outputMetrics, "outputMetrics", "output/broken.yaml", "output for the broken metrics")
	flag.StringVar(&outputRules, "outputRules", "output/rules.yaml", "output for the rules to fix broken metrics")
	flag.Parse()

	metrics, err := prometheus.GetPromMetrics(promInput)
	if err != nil {
		log.Errorf("getting prometheus metrics: %v", err)
		return
	}

	totalMetrics := len(metrics)
	brokenMetrics := listMetricWithWrongType(metrics)
	if len(brokenMetrics) == 0 {
		log.Printf("Total metrics: %d", totalMetrics)
		log.Printf("No broken metrics detected in %q", promInput)
		return
	}

	totalBroken := len(brokenMetrics)
	percentBroken := (float32(totalBroken) / float32(totalMetrics)) * 100

	log.Printf("Total metrics: %d | Broken metrics %d | Broken Metrics(%%): %f", totalMetrics, totalBroken, percentBroken)
	log.Printf("broken metrics detected in %q, placed in %q and the rules fixen them in %q", promInput, outputMetrics, outputRules)
	if err = outputBrokenMetrics(brokenMetrics, outputMetrics); err != nil {
		log.Errorf("error printing broken metrics: %s", err)
	}

	if err = outputFixRules(brokenMetrics, outputRules); err != nil {
		log.Errorf("error printing fixing rules: %s", err)
	}

}

func listMetricWithWrongType(metrics []prometheus.Metric) []prometheus.Metric {
	brokenMetrics := []prometheus.Metric{}
	computedBrokenMetrics := map[string]struct{}{}
	for _, m := range metrics {
		if isCounterWithWrongSuffix(m) || isGaugeWithWrongSuffix(m) {
			if _, ok := computedBrokenMetrics[m.Name]; !ok {
				brokenMetrics = append(brokenMetrics, m)
				computedBrokenMetrics[m.Name] = struct{}{}
			}
		}
	}

	return brokenMetrics
}

func isCounterWithWrongSuffix(m prometheus.Metric) bool {
	if m.MetricType == prometheus.MetricType_COUNTER {
		return !strings.HasSuffix(m.Name, counterSuffix_bucket) && !strings.HasSuffix(m.Name, counterSuffix_count) && !strings.HasSuffix(m.Name, counterSuffix_total)
	}
	return false
}

func isGaugeWithWrongSuffix(m prometheus.Metric) bool {
	if m.MetricType == prometheus.MetricType_GAUGE {
		return strings.HasSuffix(m.Name, counterSuffix_bucket) && strings.HasSuffix(m.Name, counterSuffix_count) && strings.HasSuffix(m.Name, counterSuffix_total)
	}
	return false
}

func outputFixRules(metrics []prometheus.Metric, filename string) error {
	file, err := openFile(filename)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprint(file, `write_relabel_configs:`)
	if err != nil {
		return fmt.Errorf("printing rules: %w", err)
	}

	for _, m := range metrics {
		nrMetricType := ""
		if m.MetricType == prometheus.MetricType_COUNTER {
			nrMetricType = nrCounter
		}
		if m.MetricType == prometheus.MetricType_GAUGE {
			nrMetricType = nrGauge
		}

		if _, err = fmt.Fprintf(file, genericRule, m.Name, nrMetricType); err != nil {
			return fmt.Errorf("printing rules: %w", err)

		}
	}
	return nil
}

func outputBrokenMetrics(metrics []prometheus.Metric, filename string) error {
	file, err := openFile(filename)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	for _, m := range metrics {
		if _, err = fmt.Fprintf(file, "%s %s\n", m.MetricType, m.Name); err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}
	}
	return nil
}

func openFile(filename string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return nil, fmt.Errorf("creating directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return file, nil
}
