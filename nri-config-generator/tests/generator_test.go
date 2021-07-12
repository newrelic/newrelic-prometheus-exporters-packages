// +build integration

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/newrelic/infra-integrations-sdk/v4/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	testIntegrationVersion = "test-tag"
	testIntegration        = "nri-powerdns"
	exporterPort           = "9120"
)

const configPDNSTemplate = `
{
  "config_protocol_version": "1",
  "action": "register_config",
  "config_name": "cfg-nri-powerdns",
  "config": {
    "variables": {},
    "integrations": [
      {
        "name": "nri-prometheus",
		{{ or .pomiInterval "" }}
        "config": {
          "standalone": false,
         {{ or .pomiVerbose "" }}
          "integration_metadata":{
            "name": "nri-powerdns",
            "version": "test-tag"
          },
          "entity_definitions": [
			{
				"conditions": [{
					"attribute":"metricName",
					"prefix":"powerdns_authoritative_"
				}],
				"identifier":"targetName", 
				"name":"targetName",
				"tags": {
					"clusterName":null,
					"targetName":null
				},
				"type":"POWERDNS_AUTHORITATIVE"
			},
			{
				"conditions": [{
					"attribute":"metricName",
					"prefix":"powerdns_recursor_"
				}],
				"identifier":"targetName", 
				"name":"targetName",
				"tags": {
					"clusterName":null,
					"targetName":null
				},
				"type":"POWERDNS_RECURSOR"
			}
          ],
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
        "name": "nri-powerdns-exporter",
        "timeout": 0,
        "exec": [
          "/usr/local/prometheus-exporters/bin/nri-powerdns-exporter",
          "--api-url",
          "http://powerdns:8080/api/v1/",
          "--api-key",
          "11111-222222-33333-44444",
          "--listen-address",
          ":{{.exporterPort}}",
          "--metric-path",
          "/metrics"
        ]
      }
    ]
  }
}
`

var defaultArgs = []string{
	"--short_running",
}
var pdnsTemplate, _ = template.New("defTemplate").Parse(configPDNSTemplate)

func TestMain(m *testing.M) {
	if err := fetchDefinitions(testIntegration); err != nil {
		panic(err.Error())
	}
	if err := buildGeneratorConfig(testIntegration, testIntegrationVersion); err != nil {
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
	templateVars := getTemplateVars(exporterPort)
	expectedResponse := executeTemplate(t, pdnsTemplate, templateVars)
	envVars := getConfigGeneratorEnvVars("config")
	stdout, err := callGeneratorConfig(testIntegration, defaultArgs, envVars)
	assert.Nil(t, err)
	assert.NotEmpty(t, stdout)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The default port is already in use, the config generator must find and available one and set the config
*/
func TestGeneratorConfigPortAlreadyInUse(t *testing.T) {
	server := &http.Server{Addr: fmt.Sprintf(":%s", exporterPort)}
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
	envVars := getConfigGeneratorEnvVars("config")
	stdout, _ := callGeneratorConfig(testIntegration, defaultArgs, envVars)
	assert.NotEmpty(t, stdout)
	assignedPort, err := getAssignedPortToPowerDNSIntegration(stdout)
	assert.Nil(t, err)
	vars := map[string]string{
		"exporterPort": assignedPort,
	}
	expectedResponse := executeTemplate(t, pdnsTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The env var interval is provided
*/
func TestGeneratorConfigWithInterval(t *testing.T) {
	envVars := getConfigGeneratorEnvVars("config")
	envVars = append(envVars, "interval=10s")
	stdout, _ := callGeneratorConfig(testIntegration, defaultArgs, envVars)
	assert.NotEmpty(t, stdout)
	vars := map[string]string{
		"exporterPort": exporterPort,
		"pomiInterval": "\"interval\":\"10s\",",
	}
	expectedResponse := executeTemplate(t, pdnsTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

/**
The verbose mode of the Agent gets propagated to exporter and prometheus
*/
func TestGeneratorVerboseMode(t *testing.T) {
	envVars := getConfigGeneratorEnvVars("config")
	envVars = append(envVars, "VERBOSE=1")
	stdout, _ := callGeneratorConfig(testIntegration, defaultArgs, envVars)
	assert.NotEmpty(t, stdout)
	vars := map[string]string{
		"exporterPort": exporterPort,
		"pomiVerbose":  "\"verbose\":\"1\",",
	}
	expectedResponse := executeTemplate(t, pdnsTemplate, vars)
	assert.JSONEq(t, expectedResponse, string(stdout))
}

func getAssignedPortToPowerDNSIntegration(content []byte) (string, error) {
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
	args, ok := integrations[1].(map[string]interface{})["exec"].([]interface{})
	if !ok {
		return "", errors.New("missing attribute config/integrations[1].env")
	}

	for i, arg := range args {
		if fmt.Sprintf("%v", arg) == "--listen-address" {
			port := strings.Trim(fmt.Sprintf("%v", args[i+1]), ":")
			return port, nil
		}
	}

	return "", fmt.Errorf("missing attribute port")
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
