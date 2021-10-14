package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetExporterNameFromIntegration(t *testing.T) {
	integrationName := "nri-powerdns"
	exporterName := getExporterNameFromIntegration(integrationName)
	assert.Equal(t, "powerdns-exporter", exporterName)
}
