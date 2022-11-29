package main

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	templateExtension = ".tmpl"
)

var (
	//go:embed integration-tests/testdata/templates
	TestTemplates embed.FS
)

func TestGetExporterNameFromIntegration(t *testing.T) {
	integrationName := "nri-powerdns"
	exporterName := getExporterNameFromIntegration(integrationName)
	assert.Equal(t, "powerdns-exporter", exporterName)
}

func Test_getExporterConfigFiles(t *testing.T) {
	expected := []string{
		"config.toml.tmpl",
	}

	got, err := getExporterConfigFiles(TestTemplates, "integration-tests/testdata/templates/exporter-config-files")

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func Test_getExporterConfigFilesEmptyFolder(t *testing.T) {
	expected := []string(nil)

	got, err := getExporterConfigFiles(TestTemplates, "integration-tests/testdata/templates/exporter-config-files-empty")

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func Test_getExporterConfigFilesNonExistentFolder(t *testing.T) {
	expected := []string(nil)

	got, err := getExporterConfigFiles(TestTemplates, "integration-tests/testdata/templates/non-existent-folder")

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func Test_generateExporterConfigFile(t *testing.T) {
	tempPath := t.TempDir()
	testTemplateFile := "config.toml.tmpl"
	testConfigFilesPath := "integration-tests/testdata/templates/exporter-config-files"
	testFile := path.Join(testConfigFilesPath, testTemplateFile)
	testOutputFile := filepath.Join(tempPath, strings.TrimSuffix(testTemplateFile, templateExtension))
	testVars := map[string]interface{}{
		"exporter_port": "9120",
	}

	err := generateExporterConfigFile(TestTemplates, testFile, tempPath, testVars)

	require.NoError(t, err)
	assert.FileExists(t, testOutputFile)
}
