// +build integration

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/newrelic/infra-integrations-sdk/v4/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const configRavenDBTemplate = `
{
  "config_protocol_version": "1",
  "action": "register_config",
  "config_name": "cfg-ravendb",
  "config": {
    "variables": {},
    "integrations": [
      {
        "name": "nri-prometheus",
		{{ or .pomiInterval "" }}
        "config": {
          "standalone": false,
		  {{ or .pomiVerbose "" }}
          "targets": [
            {
              "urls": [
                "http://localhost:{{ .exporterPort }}"
              ]
            }
          ]
        }
      },
      {
        "name": "prometheus-exporter-ravendb",
        "exec": [
          "/usr/local/prometheus-exporters/bin/ravendb-exporter"
        ],
        "timeout": 0,
        "env": {
		  {{ or .ravenVerbose "" }}
          "RAVENDB_URL": "http://live-test.ravendb.net",
          "PORT": "{{ .exporterPort }}"
        }
      }
    ]
  }
}
`

var ravendbTemplate, _ = template.New("defTemplate").Parse(configRavenDBTemplate)

func TestMain(m *testing.M) {
	if err := buildGeneratorConfig("ravendb"); err != nil {
		panic(err.Error())
	}
	exitVal := m.Run()
	clean()
	os.Exit(exitVal)
}

/**
Happy path
*/
func TestGeneratorConfig(t *testing.T) {
	integration := "ravendb"
	vars := map[string]string{
		"exporterPort": "3333",
	}
	expectedResponse := executeTemplate(t, ravendbTemplate, vars)
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, _ := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.NotEmpty(t, stdout)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The default port is already in use, the config generator must find and available one and set the config
*/
func TestGeneratorConfigPortAlreadyInUse(t *testing.T) {
	integration := "ravendb"
	defaultPort := "3333"
	server := &http.Server{Addr: fmt.Sprintf(":%s", defaultPort)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Warn(err.Error())
		}
	}()
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, _ := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.NotEmpty(t, stdout)
	assignedPort, err := getAssignedPortToExporter(stdout)
	assert.Nil(t, err)
	vars := map[string]string{
		"exporterPort": assignedPort,
	}
	expectedResponse := executeTemplate(t, ravendbTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The env var interval is provided
*/
func TestGeneratorConfigWithInterval(t *testing.T) {
	integration := "ravendb"

	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	intervalEnvVar := "interval=10s"
	stdout, _ := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar, intervalEnvVar})
	assert.NotEmpty(t, stdout)
	assignedPort, err := getAssignedPortToExporter(stdout)
	assert.Nil(t, err)
	vars := map[string]string{
		"exporterPort": assignedPort,
		"pomiInterval": "\"interval\":\"10s\",",
	}
	expectedResponse := executeTemplate(t, ravendbTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The verbose mode of the Agent gets propagated to exporter and prometheus
*/
func TestGeneratorVerboseMode(t *testing.T) {
	integration := "ravendb"

	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	verboseEnvVar := "VERBOSE=1"
	stdout, _ := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar, verboseEnvVar})
	assert.NotEmpty(t, stdout)
	assignedPort, err := getAssignedPortToExporter(stdout)
	assert.Nil(t, err)
	vars := map[string]string{
		"exporterPort": assignedPort,
		"ravenVerbose": "\"VERBOSE\":\"1\",",
		"pomiVerbose":  "\"verbose\":\"1\",",
	}
	expectedResponse := executeTemplate(t, ravendbTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The export port is defined by the user
*/
func TestGeneratorConfigWithExporterPortInConfigFile(t *testing.T) {
	integration := "ravendb"
	vars := map[string]string{
		"exporterPort": "3333",
	}
	expectedResponse := executeTemplate(t, ravendbTemplate, vars)
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config_with_port.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, _ := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.NotEmpty(t, stdout)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The export port is defined by the user but it's already in use
*/
func TestGeneratorConfigWithExporterPortInConfigFileButItsInUse(t *testing.T) {
	integration := "ravendb"

	server := &http.Server{Addr: fmt.Sprintf(":%s", "3333")}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Warn(err.Error())
		}
	}()
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	configPath := filepath.Join(rootDir(), "tests", "testdata", "config_with_port.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, _ := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.NotEmpty(t, stdout)
	out := make(map[string]interface{})
	json.Unmarshal(stdout, &out)
	cfg := out["config"].(map[string]interface{})
	integrations := cfg["integrations"].([]interface{})
	env := integrations[1].(map[string]interface{})["env"].(map[string]interface{})
	port := env["PORT"]
	vars := map[string]string{
		"exporterPort": fmt.Sprintf("%s", port),
	}
	expectedResponse := executeTemplate(t, ravendbTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

func getAssignedPortToExporter(content []byte) (string, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(content, &v); err != nil {
		return "", err
	}
	config, ok := v["config"].(map[string]interface{})
	if !ok {
		return "", errors.New("missing attribute config")
	}
	integrations, ok := config["integrations"].([]interface{})
	if !ok {
		return "", errors.New("missing attribute config/integrations")
	}
	if len(integrations) != 2 {
		return "", errors.New("missing attribute config/integrations[1]")
	}
	env, ok := integrations[1].(map[string]interface{})["env"].(map[string]interface{})
	if !ok {
		return "", errors.New("missing attribute config/integrations[1].env")
	}
	return fmt.Sprintf("%v", env["PORT"]), nil
}

func executeTemplate(t *testing.T, tpl *template.Template, vars map[string]string) string {
	if tpl == nil {
		t.Fatal("invalid template")
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, vars); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
