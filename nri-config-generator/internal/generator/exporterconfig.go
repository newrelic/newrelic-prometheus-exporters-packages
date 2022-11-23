package generator

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type ExporterConfig interface {
	Generate(vars map[string]interface{}) (string, error)
}

type exporterConfig struct {
	template *template.Template
}

func NewExporterConfig(template *template.Template) ExporterConfig {
	return &exporterConfig{
		template: template,
	}
}

func (g *exporterConfig) Generate(vars map[string]interface{}) (string, error) {
	var templateOut bytes.Buffer
	if err := g.template.Execute(&templateOut, vars); err != nil {
		log.Errorf("error executing the template for the exporter config: '%s'", err.Error())
		return "", err
	}
	return templateOut.String(), nil
}
