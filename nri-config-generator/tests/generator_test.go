package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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

func callGeneratorConfig(integration string, args []string, env []string) ([]byte, error) {
	executable := fmt.Sprintf("nri-%s", integration)
	cmd := &exec.Cmd{
		Path: filepath.Join(rootDir(), "bin", executable),
		Args: args,
		Env:  env,
	}
	return cmd.Output()
}

/**
Happy path
*/
func TestGeneratorConfig(t *testing.T) {
	integration := "ravendb"
	defaultPort := "3333"
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", defaultPort, defaultPort)
	assert.Nil(t, buildGeneratorConfig(integration, defaultPort))
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, err := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.Nil(t, err)
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
	assert.Nil(t, buildGeneratorConfig(integration, defaultPort))
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, err := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.Nil(t, err)
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
	defaultPort := "3333"
	assert.Nil(t, buildGeneratorConfig(integration, defaultPort))
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	intervalEnvVar := "interval=10s"
	stdout, err := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar, intervalEnvVar})
	assert.Nil(t, err)
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
	defaultPort := "3333"
	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", "9911", "9911")
	assert.Nil(t, buildGeneratorConfig(integration, defaultPort))
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config_with_port.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, err := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.Nil(t, err)
	assert.NotEmpty(t, stdout)
	assert.JSONEq(t, expectedResponse, string(stdout))
}


/**
The export port is defined by the user but it's already in use
*/
func TestGeneratorConfigWithExporterPortInConfigFileButItsInUse(t *testing.T) {
	integration := "ravendb"
	defaultPort := "3333"

	server := &http.Server{Addr: fmt.Sprintf(":%s", "9911")}
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

	expectedResponse := fmt.Sprintf(configRavenDBTemplate, "", "3333", "3333")
	assert.Nil(t, buildGeneratorConfig(integration, defaultPort))
	configPath := filepath.Join(rootDir(), "tests", "testdata", "config_with_port.yml")
	configPathEnvVar := fmt.Sprintf("CONFIG_PATH=%s", configPath)
	stdout, err := callGeneratorConfig(integration, []string{}, []string{configPathEnvVar})
	assert.Nil(t, err)
	assert.NotEmpty(t, stdout)
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
