package main

import (
	"embed"
	"fmt"
	"os"
	"runtime"
	"text/template"
	"time"

	"github.com/newrelic/infra-integrations-sdk/v4/log"
	"github.com/newrelic/nri-config-generator/args"
	"github.com/newrelic/nri-config-generator/generator"
	"github.com/newrelic/nri-config-generator/httport"
	"github.com/newrelic/nri-config-generator/synthesis"
)

const (
	varIntegrationName      = "integration"
	varExporterPort         = "exporter_port"
	varSynthesisDefinitions = "entity_definitions"
	varExporterDefinition   = "exporter_definition"
	sleepTime               = 30 * time.Second
	definitionFileName      = "definitions/definitions.yml"
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
	//go:embed definitions
	definitions embed.FS
	vars        = map[string]interface{}{
		varIntegrationName: integration,
	}
)

func main() {
	al, err := args.PopulateVars(vars)
	panicErr(err)
	if al.ShowVersion {
		printVersion()
		os.Exit(0)
	}
	integrationTemplatePattern := fmt.Sprintf("templates/%s.json.tmpl", integration)
	content, err := integrationTemplate.ReadFile(integrationTemplatePattern)
	panicErr(err)
	integrationTemplate, err := loadIntegrationTemplate(content)
	panicErr(err)
	exporterGenerator := generator.NewExporter(integration, integrationTemplate)
	configTemplate, err := loadConfigTemplate()
	panicErr(err)
	configGenerator := generator.NewConfig(configTemplate)
	port, err := findExporterPort()
	panicErr(err)
	definitionContent, err := definitions.ReadFile(definitionFileName)
	panicErr(err)
	definitions, err := synthesis.ProcessSynthesisDefinitions(definitionContent)
	panicErr(err)
	vars[varSynthesisDefinitions] = definitions
	vars[varExporterPort] = fmt.Sprintf("%v", port)
	exporterText, err := exporterGenerator.Generate(vars)
	panicErr(err)
	vars[varExporterDefinition] = exporterText
	output, err := configGenerator.Generate(vars)
	panicErr(err)
	fmt.Println(output)
	// In this case the integration is completed after fetching once the data from the exporter
	if al.ShortRunning {
		os.Exit(0)
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

func panicErr(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf(
		"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
		integration,
		integrationVersion,
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		runtime.Version(),
		gitCommit,
		buildDate)
}

func findExporterPort() (int, error) {
	cfgPort := ""
	if cfg, ok := vars[args.PrefixCfg]; ok {
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
