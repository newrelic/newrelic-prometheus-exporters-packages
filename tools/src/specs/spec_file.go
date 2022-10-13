package specs

import (
	"fmt"
	"github.com/newrelic/newrelic-prometheus-exporters-packages/tools/src/prometheus"
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	internalAttributes = []string{"newrelic.integrationName", "newrelic.integrationVersion", "newrelic.source", "newrelic.agentVersion"}
	ignoredAttributes  = []string{"agentName", "coreCount", "processorCount", "systemMemoryBytes"}
)

type Specs struct {
	SpecVersion                  string    `yaml:"specVersion"`
	OwningTeam                   string    `yaml:"owningTeam"`
	IntegrationName              string    `yaml:"integrationName"`
	HumanReadableIntegrationName string    `yaml:"humanReadableIntegrationName"`
	Entities                     []*Entity `yaml:"entities"`
}

// Metrics
type MetricSpec struct {
	Name              string      `yaml:"name"`
	Type              string      `yaml:"type"`
	DefaultResolution int         `yaml:"defaultResolution"`
	Unit              string      `yaml:"unit"`
	Description       string      `yaml:"-"`
	Dimensions        []Dimension `yaml:"dimensions,omitempty"`
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
	Tags               []string            `yaml:"tags"`
}

type Dimension struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func GenerateSpecFile(c Config, metrics []prometheus.Metric, filename string) *Specs {
	sp := Specs{
		SpecVersion:                  "2",
		OwningTeam:                   "integrations",
		IntegrationName:              c.IntegrationName,
		HumanReadableIntegrationName: c.HumanReadableIntegrationName,
		Entities:                     nil,
	}
	for _, m := range metrics {
		entityType, ok := getMatchingEntity(c.Definitions, m.Name, m.Attributes)
		if !ok {
			fmt.Printf("metric not matching %s \n", m.Name)
			continue
		}

		e, ok := isEntityDefined(sp.Entities, entityType)
		if !ok {
			e = &Entity{
				EntityType: entityType,
			}
			sp.Entities = append(sp.Entities, e)
		}

		if ok := isMetricDefined(e.Metrics, m.Name); ok {
			continue
		}

		mSpec := MetricSpec{
			Name:              m.Name,
			Type:              string(m.MetricType),
			DefaultResolution: 15,
			Unit:              "count",
			Description:       m.Description,
			Dimensions:        nil,
		}
		for k, _ := range m.Attributes {
			mSpec.Dimensions = append(mSpec.Dimensions, Dimension{
				Name: k,
				Type: "string",
			})
		}

		if mSpec.Type == string(prometheus.MetricType_HISTOGRAM) {
			e.Metrics = append(e.Metrics, computeHistogramMetrics(mSpec)...)
		} else {
			e.Metrics = append(e.Metrics, &mSpec)
		}
	}

	for _, e := range sp.Entities {
		for _, i := range internalAttributes {
			e.InternalAttributes = append(e.InternalAttributes, InternalAttribute{
				Name: i,
				Type: "string",
			})
		}
		e.IgnoredAttributes = ignoredAttributes
	}

	out, err := yaml.Marshal(sp)
	if err != nil {
		log.Print(err)
	}
	err = os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		log.Print(err)
	}
	err = ioutil.WriteFile(filename, out, 0777)
	if err != nil {
		log.Print(err)
	}
	return &sp
}

func computeHistogramMetrics(m MetricSpec) []*MetricSpec {
	sumDimensions := make([]Dimension, len(m.Dimensions))
	copy(sumDimensions, m.Dimensions)

	bucketDimensions := make([]Dimension, len(m.Dimensions))
	copy(bucketDimensions, m.Dimensions)
	bucketDimensions = append(bucketDimensions, Dimension{
		Name: "le",
		Type: "string",
	})

	return []*MetricSpec{
		{
			Name:              m.Name + "_sum",
			Type:              string(prometheus.MetricType_SUMMARY),
			DefaultResolution: m.DefaultResolution,
			Unit:              "count",
			Description:       m.Description + " (sum metric)",
			Dimensions:        sumDimensions,
		},
		{
			Name:              m.Name + "_bucket",
			Type:              string(prometheus.MetricType_COUNTER),
			DefaultResolution: m.DefaultResolution,
			Unit:              "count",
			Description:       m.Description + " (bucket metric)",
			Dimensions:        bucketDimensions,
		},
	}
}

// getMatchingRule iterates over all conditions to check if m satisfy returning the associated rule.
func getMatchingEntity(definitions []Definition, metricName string, attributes prometheus.Set) (string, bool) {
	matchedCondition := &Condition{}
	matchedEntityName := ""

	for _, d := range definitions {
		for _, c := range d.Conditions {
			// special case since metricName is not a metric attribute.
			value := metricName

			if c.Attribute != "metricName" {
				val, ok := attributes[c.Attribute]
				if !ok {
					continue
				}

				value, _ = val.(string)
			}

			// longer prefix matches take precedences over shorter ones.
			// this allows to discriminate "foo_bar_" from "foo_" kind of metrics.
			if c.Match(value) && (matchedEntityName == "" || len(c.Prefix) > len(matchedCondition.Prefix)) { // nosemgrep: bad-nil-guard
				matchedCondition = &c
				matchedEntityName = d.EntityType
			}
		}
	}

	return matchedEntityName, matchedEntityName != ""
}

func isEntityDefined(entities []*Entity, entityType string) (*Entity, bool) {
	for _, entity := range entities {
		if entity.EntityType == entityType {
			return entity, true
		}
	}

	return nil, false
}

func isMetricDefined(metrics []*MetricSpec, metricName string) bool {
	for _, m := range metrics {
		if m.Name == metricName {
			return true
		}
	}

	return false
}
