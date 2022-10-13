package docs

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/newrelic/newrelic-prometheus-exporters-packages/tools/src/specs"
)

const (
	funcSet = "set"
)

//go:embed input/template.tmpl
var docTemplateContent string

func loadDocTemplate() *template.Template {
	t, err := template.New("").Funcs(TemplatesFunc).Parse(docTemplateContent)
	if err != nil {
		panic(err)
	}
	return t
}

var TemplatesFunc = template.FuncMap{
	funcSet: set,
}

func set(property string, value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("\"%s\":%s,", property, attributeValue(value))
}

func attributeValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func GenerateDocFile(sp *specs.Specs, fileName string) {
	template := loadDocTemplate()
	r, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer r.Close()

	err = template.Execute(r, sp)
	if err != nil {
		fmt.Printf("executing template: %s", err.Error())
	}
}
