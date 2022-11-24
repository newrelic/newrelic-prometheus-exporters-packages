package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"text/template"

	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
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
	exporterConfigFilesPath = "templates/exporter-config-files"
	nixExportsBinPath       = "/usr/local/prometheus-exporters/bin"
	winExportsBinPath       = "C:\\Program Files\\Prometheus-exporters\\bin"
	sleepTime               = 30 * time.Second
	emptyMap                = "[]"
	templateSuffix          = ".tmpl"
)

var (
	integration        string
	integrationVersion string
	gitCommit          string
	buildDate          string
	//go:embed templates
	exporterConfigTemplates embed.FS
	//go:embed templates
	integrationTemplate embed.FS
	//go:embed templates
	configTemplate embed.FS
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

	vars, exporterConfigPath, err := config.GetVars(al)
	integration = "nri-powerdns"
	if err != nil {
		log.Fatal(err)
	}
	exporterName := getExporterNameFromIntegration(integration)

	exporterConfigFiles, err := getExporterConfigFiles()
	if err != nil {
		log.Fatal(err)
	}

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

	// Create exporter config files
	for _, configFile := range exporterConfigFiles {
		err = generateExporterConfigFile(configFile, exporterConfigPath, vars)
		if err != nil {
			log.Fatal(err)
		}
	}

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
	// long-running execution. This is a long-running execution because in case of the export port was auto-generated (a random port inc ase
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

func getExporterConfigFiles() ([]string, error) {
	var configFiles []string
	fileList, err := fs.ReadDir(exporterConfigTemplates, exporterConfigFilesPath)
	if err != nil {
		return nil, err
	}

	for _, file := range fileList {
		configFiles = append(configFiles, file.Name())
	}
	return configFiles, nil
}

func generateExporterConfigFile(exporterTemplateFile string, exporterConfigPath string, vars map[string]interface{}) error {
	templateLocation := fmt.Sprintf("%s/%s", exporterConfigFilesPath, exporterTemplateFile)

	content, err := exporterConfigTemplates.ReadFile(templateLocation)
	if err != nil {
		return fmt.Errorf("reading exporter config template %s, %w", exporterTemplateFile, err)
	}

	exporterConfigTemplate, err := loadTemplate("exporter config", content)
	if err != nil {
		return fmt.Errorf("loadConfigTemplate, %w", err)
	}

	exporterConfigGenerator := generator.NewExporterConfig(exporterConfigTemplate)
	result, err := exporterConfigGenerator.Generate(vars)
	if err != nil {
		return fmt.Errorf("exporterConfigGenerator.Generate: %w", err)
	}

	filename := strings.TrimSuffix(exporterTemplateFile, templateSuffix)
	path := sanitizePath(exporterConfigPath)
	outputFile := fmt.Sprintf("%s%s", path, filename)

	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("exporterConfigGenerator.Writing: %w", err)
	}

	return nil
}

func getConfigGenerator() (generator.Config, error) {
	configTemplatePattern := fmt.Sprintf("templates/%s.prometheus.json.tmpl", integration)
	content, err := configTemplate.ReadFile(configTemplatePattern)
	if err != nil {
		return nil, fmt.Errorf("readingfile %s, %w", configTemplatePattern, err)
	}

	configTemplate, err := loadTemplate("prometheus configuration", content)
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

	integrationTemplate, err := loadTemplate("integration", content)
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

func loadTemplate(templateType string, content []byte) (*template.Template, error) {
	t, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(string(content))
	if err != nil {
		log.Error("error parsing the %s template: '%s'", templateType, err.Error())
		return nil, err
	}
	return t, nil
}

func sanitizePath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%s/", path)
	}
	return path
}

func prometheusExportersBinPath(name string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(filepath.Join(winExportsBinPath, fmt.Sprintf("%s.exe", name)), "\\", "\\\\", -1)
	}
	return filepath.Join(nixExportsBinPath, name)
}
