package main

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/newrelic/infra-integrations-sdk/v4/log"
	"github.com/newrelic/nri-config-generator/args"
	"github.com/newrelic/nri-config-generator/httport"
)

const (
	varIntegrationName    = "integration"
	varExporterPort       = "exporter_port"
	varExporterDefinition = "exporter_definition"
)

var (
	integration string
	defPort     string
	//go:embed templates
	integrationTemplate embed.FS
	//go:embed templates/config.json.tmpl
	configTemplateContent string
	configTemplate        = template.New("")
	vars                  = map[string]interface{}{
		varIntegrationName: integration,
	}
)

func main() {
	if err := args.PopulateVars(vars); err != nil {
		log.Error("error populating vars: '%s'", err.Error())
		panic(err)
	}
	cfgPort := ""
	if cfg, ok := vars[args.PrefixCfg]; ok {
		cfgVars := cfg.(map[string]interface{})
		if cfgPortPtr := cfgVars[varExporterPort]; cfgPortPtr != nil {
			cfgPort = fmt.Sprintf("%v", cfgPortPtr)
		}
	}
	port, err := httport.GetAvailablePort(cfgPort, defPort)
	if err != nil {
		log.Error("error obtaining the port for the exporter: '%s'", err.Error())
		panic(err)
	}

	vars[varExporterPort] = port
	integrationTemplatePattern := fmt.Sprintf("templates/%s.json.tmpl", integration)

	t, err := template.ParseFS(integrationTemplate, integrationTemplatePattern)
	if err != nil {
		log.Error("error parsing the integration template: '%s'", err.Error())
		panic(err)
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, vars); err != nil {
		log.Error("error executing the template for the integration: '%s'", err.Error())
		panic(err)
	}
	vars[varExporterDefinition] = tpl.String()
	t, err = configTemplate.Parse(configTemplateContent)
	if err != nil {
		log.Error("error parsing the template for the config: '%s'", err.Error())
		panic(err)
	}
	var output bytes.Buffer
	if err := t.Execute(&output, vars); err != nil {
		log.Error("error executing the template for the config: '%s'", err.Error())
		panic(err)
	}
	fmt.Println(strings.ReplaceAll(output.String(), "\n", ""))
}
