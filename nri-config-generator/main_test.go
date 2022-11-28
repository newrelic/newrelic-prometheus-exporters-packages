package main

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"os"
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

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

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

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func Test_getExporterConfigFilesEmptyFolder(t *testing.T) {
	expected := []string(nil)

	got, err := getExporterConfigFiles(TestTemplates, "integration-tests/testdata/templates/exporter-config-files-empty")

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func Test_getExporterConfigFilesNonExistentFolder(t *testing.T) {
	expected := []string(nil)

	got, err := getExporterConfigFiles(TestTemplates, "integration-tests/testdata/templates/non-existent-folder")

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func Test_generateExporterConfigFile(t *testing.T) {
	tempPath := os.TempDir()
	templateFile := "config.toml.tmpl"
	testFile := filepath.Join(tempPath, strings.TrimSuffix(templateFile, templateExtension))
	testVars := map[string]interface{}{
		"exporter_port": "9120",
	}
	err := generateExporterConfigFile(templateFile, tempPath, testVars)
	if err != nil {
		t.Error(err)
	}
	assert.FileExists(t, testFile)
}
