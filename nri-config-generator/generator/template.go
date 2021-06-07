package generator

import (
	"fmt"
	"text/template"
)

const (
	funcSet = "set"
)

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
