package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	internalAttributes = []string{"newrelic.integrationName", "newrelic.integrationVersion", "newrelic.source", "newrelic.agentVersion",
		"instrumentation.name", "instrumentation.provider", "instrumentation.version", "collector.name", "collector.version"}
	ignoredAttributes = []string{"agentName", "coreCount", "processorCount", "systemMemoryBytes"}
)

func generateSpecFile(c Config, metrics []Metric, filename string) *Specs {
	sp := Specs{
		SpecVersion:                  "2",
		OwningTeam:                   "integrations",
		IntegrationName:              c.IntegrationName,
		HumanReadableIntegrationName: c.HumanReadableIntegrationName,
		Entities:                     nil,
	}
	for _, m := range metrics {
		entityType, ok := getMatchingEntity(c.Definitions, m.name, m.attributes)
		if !ok {
			fmt.Printf("metric not matching %s \n", m.name)
			continue
		}

		e, ok := isEntityDefined(sp.Entities, entityType)
		if !ok {
			e = &Entity{
				EntityType: entityType,
			}
			sp.Entities = append(sp.Entities, e)
		}

		if ok := isMetricDefined(e.Metrics, m.name); ok {
			continue
		}

		mSpec := MetricSpec{
			Name:              m.name,
			Type:              string(m.metricType),
			DefaultResolution: 15,
			Unit:              "[TO_BE_FILLED with one of count, percent, seconds, milliseconds, microsenconds, bytes, countPerSecond, bytesPerSecond, bitPerSecond]",
			Description:       m.description,
			Dimensions:        nil,
		}
		for k, _ := range m.attributes {
			mSpec.Dimensions = append(mSpec.Dimensions, Dimension{
				Name: k,
				Type: "string",
			})
		}
		e.Metrics = append(e.Metrics, &mSpec)
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
	err = ioutil.WriteFile(filename, out, 0666)
	if err != nil {
		log.Print(err)
	}
	return &sp
}

// getMatchingRule iterates over all conditions to check if m satisfy returning the associated rule.
func getMatchingEntity(definitions []Definition, metricName string, attributes Set) (string, bool) {
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
			if c.match(value) && (matchedEntityName == "" || len(c.Prefix) > len(matchedCondition.Prefix)) { // nosemgrep: bad-nil-guard
				matchedCondition = &c
				matchedEntityName = d.EntityType
			}
		}
	}

	return matchedEntityName, matchedEntityName != ""
}

// match evaluates the condition an attribute by checking either its whole name against `Value` or if it starts with `Prefix`.
func (c Condition) match(attribute string) bool {
	if c.Value != "" {
		return c.Value == attribute
	}
	// this returns true if c.Prefix is "" and is ok since the attribute exists
	return strings.HasPrefix(attribute, c.Prefix)
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
