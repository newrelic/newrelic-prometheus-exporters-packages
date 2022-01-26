package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	sdkv4 "github.com/newrelic/infra-integrations-sdk/v4/args"
	"github.com/newrelic/infra-integrations-sdk/v4/log"

	"github.com/newrelic/nri-config-generator/internal/config"
	"github.com/newrelic/nri-config-generator/internal/generator"
	"github.com/newrelic/nri-config-generator/internal/httport"
)

const (
	varIntegrationName      = "integration"
	varIntegrationVersion   = "integration_version"
	varExporterPort         = "exporter_port"
	varExporterDefinition   = "exporter_definition"
	varExporterBinaryPath   = "exporter_binary_path"
	varMetricTransformation = "transformations"
	nixExportsBinPath       = "/usr/local/prometheus-exporters/bin"
	winExportsBinPath       = "C:\\Program Files\\Prometheus-exporters\\bin"
	sleepTime               = 30 * time.Second
	emptyMap                = "[]"
)

var (
	integration        string
	integrationVersion string
	gitCommit          string
	buildDate          string
	//go:embed templates
	integrationTemplate embed.FS
	//go:embed templates/config.json.tmpl
	configTemplateContent string
)

func main() {

	al := &config.ArgumentList{}
	err := sdkv4.SetupArgs(al)
	if err != nil {
		log.Fatal(err)
	}

	if al.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			integration, integrationVersion, fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(), gitCommit, buildDate)
		return
	}

	vars, err := config.GetVars(al)
	if err != nil {
		log.Fatal(err)
	}
	exporterName := getExporterNameFromIntegration(integration)
	exporterGenerator, err := getExporterGenerator(exporterName)
	if err != nil {
		log.Fatal(err)
	}

	configGenerator, err := getConfigGenerator()
	if err != nil {
		log.Fatal(err)
	}

	transformations, err := getMetricTransformations(vars)
	if err != nil {
		log.Fatal(err)
	}

	port, err := findExporterPort(vars)
	if err != nil {
		log.Fatal(err)
	}

	vars[varExporterBinaryPath] = prometheusExportersBinPath(exporterName)
	vars[varExporterPort] = fmt.Sprintf("%v", port)
	vars[varMetricTransformation] = transformations
	vars[varIntegrationName] = integration
	vars[varIntegrationVersion] = integrationVersion

	output, err := generateOutput(exporterGenerator, configGenerator, vars)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
	// In this case the integration is completed after fetching once the data from the exporter
	if al.ShortRunning {
		return
	}

	httport.SetPrometheusExporterPort("localhost", port)
	// long running execution. This is a long running execution because in case of the export port was auto-generated (a random port inc ase
	// of the exporter_port is not provided by configuration) we need to keep this port for the exporter.
	for {
		time.Sleep(sleepTime)
		fmt.Println("{}")
		fmt.Println(output)
	}
}

func generateOutput(exporterGenerator generator.Exporter, configGenerator generator.Config, vars map[string]interface{}) (string, error) {
	exporterText, err := exporterGenerator.Generate(vars)
	if err != nil {
		return "", fmt.Errorf("exporterGenerator.Generate: %w", err)
	}

	vars[varExporterDefinition] = exporterText

	output, err := configGenerator.Generate(vars)
	if err != nil {
		return "", fmt.Errorf("configGenerator.Generate: %w", err)
	}

	return output, err
}

func getExporterNameFromIntegration(integration string) string {
	var exporterNameBeforeFix = fmt.Sprintf("%s-exporter", integration)
	return strings.Replace(exporterNameBeforeFix, "nri-", "", 1)
}

func getConfigGenerator() (generator.Config, error) {
	configTemplate, err := loadConfigTemplate()
	if err != nil {
		return nil, fmt.Errorf("loadConfigTemplate, %w", err)
	}
	configGenerator := generator.NewConfig(configTemplate)
	return configGenerator, nil
}

func getExporterGenerator(exporterName string) (generator.Exporter, error) {
	integrationTemplatePattern := fmt.Sprintf("templates/%s.json.tmpl", integration)
	content, err := integrationTemplate.ReadFile(integrationTemplatePattern)
	if err != nil {
		return nil, fmt.Errorf("readingfile %s, %w", integrationTemplatePattern, err)
	}

	integrationTemplate, err := loadIntegrationTemplate(content)
	if err != nil {
		return nil, fmt.Errorf("loadIntegrationTemplate, %w", err)
	}
	exporterGenerator := generator.NewExporter(exporterName, integrationTemplate)
	return exporterGenerator, nil
}

func getMetricTransformations(vars map[string]interface{}) (string, error) {
	if cfg, ok := vars[config.PrefixCfg]; ok {
		cfgVars := cfg.(map[string]interface{})
		if metricTransformations, ok := cfgVars[varMetricTransformation]; ok {

			mT := &[]config.ProcessingRule{}
			err := mapstructure.Decode(metricTransformations, mT)
			if err != nil {
				return "", fmt.Errorf("mapstructure decoding: '%v', %w", metricTransformations, err)
			}
			pBytes, err := json.Marshal(*mT)
			if err != nil {
				return "", fmt.Errorf("marshalling: '%v', %w", mT, err)
			}
			return string(pBytes), nil
		}
	}
	return emptyMap, nil
}

func findExporterPort(vars map[string]interface{}) (int, error) {
	cfgPort := ""
	if cfg, ok := vars[config.PrefixCfg]; ok {
		cfgVars := cfg.(map[string]interface{})
		if cfgPortPtr := cfgVars[varExporterPort]; cfgPortPtr != nil {
			cfgPort = fmt.Sprintf("%v", cfgPortPtr)
		}
	}
	port, err := httport.GetAvailablePort(cfgPort)
	if err != nil {
		log.Error("error obtaining the port for the exporter: '%s'", err.Error())
		return -1, err
	}
	return port, nil
}

func loadIntegrationTemplate(content []byte) (*template.Template, error) {
	t, err := template.New("").Funcs(generator.TemplatesFunc).Parse(string(content))
	if err != nil {
		log.Error("error parsing the integration template: '%s'", err.Error())
		return nil, err
	}
	return t, nil
}

func loadConfigTemplate() (*template.Template, error) {
	t, err := template.New("").Funcs(generator.TemplatesFunc).Parse(configTemplateContent)
	if err != nil {
		log.Error("error parsing the template for the config: '%s'", err.Error())
		return nil, err
	}
	return t, nil
}

func prometheusExportersBinPath(name string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(winExportsBinPath, fmt.Sprintf("%s.exe", name))
	}
	return filepath.Join(nixExportsBinPath, name)
}
