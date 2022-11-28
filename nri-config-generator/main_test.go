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
	//go:embed templates
	TestTemplates embed.FS
)

func copyTemplate(fileName string) error {
	sourcePath := filepath.Join("integration-tests", "testdata", "integration_template", fileName)
	bytesRead, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	targetPath := filepath.Join("templates", "exporter-config-files", fileName)
	err = os.WriteFile(targetPath, bytesRead, 0600)
	if err != nil {
		return err
	}
	return nil
}

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
	err := copyTemplate("config.toml.tmpl")
	if err != nil {
		t.Error(err)
	}
	expected := []string{
		"config.toml.tmpl",
	}

	got, err := getExporterConfigFiles(TestTemplates, "templates/exporter-config-files")

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
