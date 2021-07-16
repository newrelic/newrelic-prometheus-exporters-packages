package generator

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type Config interface {
	Generate(vars map[string]interface{}) (string, error)
}

type config struct {
	template *template.Template
}

func NewConfig(template *template.Template) Config {
	return &config{
		template: template,
	}
}

func (g *config) Generate(vars map[string]interface{}) (string, error) {
	var templateOut bytes.Buffer
	if err := g.template.Execute(&templateOut, vars); err != nil {
		log.Errorf("error executing the template for the config: '%s'", err.Error())
		return "", err
	}
	content := compactTextInOneLine(templateOut.String())
	return removeTrailingCommas(content), nil
}
