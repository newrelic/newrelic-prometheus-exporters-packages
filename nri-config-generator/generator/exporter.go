package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var integrationDetailsAttributes = `"name": "%s-exporter","exec": ["%s"],"timeout": 0,`

type Exporter interface {
	Generate(vars map[string]interface{}) (string, error)
}

type exporter struct {
	details  string
	template *template.Template
}

func NewExporter(name string, template *template.Template) Exporter {
	exporterExec := prometheusExportersBinPath(name)
	return &exporter{
		details:  fmt.Sprintf(integrationDetailsAttributes, name, exporterExec),
		template: template,
	}
}

func (g *exporter) Generate(vars map[string]interface{}) (string, error) {
	var templateOut bytes.Buffer
	if err := g.template.Execute(&templateOut, vars); err != nil {
		log.Error("error executing the template for the integration: '%s'", err.Error())
		return "", err
	}
	content := g.appendIntegrationDetails(templateOut.String())
	content = compactTextInOneLine(content)
	content = removeTrailingCommas(content)
	return content, nil
}

func (g *exporter) appendIntegrationDetails(content string) string {
	return strings.Replace(content, "{", fmt.Sprintf("{%s", g.details), 1)
}
