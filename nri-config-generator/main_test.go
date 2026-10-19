package main

import (
	"embed"
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_resolveExporterBinPath(t *testing.T) {
	const name = "some-exporter"

	t.Run("uses default dir when the binary exists there", func(t *testing.T) {
		defaultDir := t.TempDir()
		require.NoError(t, os.WriteFile(exporterBinPath(defaultDir, name), nil, 0o755))
		executable := func() (string, error) {
			t.Fatal("executable must not be queried when the binary exists in the default dir")
			return "", nil
		}

		got := resolveExporterBinPath(defaultDir, name, executable)

		assert.Equal(t, exporterBinPath(defaultDir, name), got)
	})

	t.Run("falls back to the executable dir when the binary is missing", func(t *testing.T) {
		defaultDir, exeDir := t.TempDir(), t.TempDir()
		executable := func() (string, error) { return filepath.Join(exeDir, "nri-integration"), nil }

		got := resolveExporterBinPath(defaultDir, name, executable)

		assert.Equal(t, exporterBinPath(exeDir, name), got)
	})

	t.Run("keeps the default path when the default dir cannot be accessed", func(t *testing.T) {
		if runtime.GOOS == "windows" || os.Geteuid() == 0 {
			t.Skip("directory permissions are not enforced on windows or for root")
		}
		defaultDir := t.TempDir()
		require.NoError(t, os.Chmod(defaultDir, 0o000))
		t.Cleanup(func() { _ = os.Chmod(defaultDir, 0o700) })
		executable := func() (string, error) {
			t.Fatal("executable must not be queried when the default dir cannot be accessed")
			return "", nil
		}

		got := resolveExporterBinPath(defaultDir, name, executable)

		assert.Equal(t, exporterBinPath(defaultDir, name), got)
	})

	t.Run("keeps the default path when the executable cannot be determined", func(t *testing.T) {
		defaultDir := t.TempDir()
		executable := func() (string, error) { return "", errors.New("boom") }

		got := resolveExporterBinPath(defaultDir, name, executable)

		assert.Equal(t, exporterBinPath(defaultDir, name), got)
	})
}
