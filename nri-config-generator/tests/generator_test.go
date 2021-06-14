// +build integration

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

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
		%s
        "config": {
          "standalone": false,
          "targets": [
            {
              "urls": [
                "http://localhost:%s"
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
          "RAVENDB_URL": "http://live-test.ravendb.net",
          "PORT": "%s"
        }
      }
    ]
  }
}
`

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
	defaultPort := "3333"
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", defaultPort, defaultPort)
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout,_:= callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
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
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", assignedPort, assignedPort)
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
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "\"interval\":\"10s\",", assignedPort, assignedPort)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The export port is defined by the user
*/
func TestGeneratorConfigWithExporterPortInConfigFile(t *testing.T) {
	integration := "ravendb"
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", "3333", "3333")
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
	out:=make(map[string]interface{})
	json.Unmarshal(stdout,&out)
	cfg:=out["config"].(map[string]interface{})
	integrations:=cfg["integrations"].([]interface{})
	env:=integrations[1].(map[string]interface{})["env"].(map[string]interface{})
	port:=env["PORT"]
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", port, port)
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
