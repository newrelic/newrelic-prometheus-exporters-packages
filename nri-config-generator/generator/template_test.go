package generator

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func expand(tmpl *template.Template, params map[string]interface{}) string {
	var output bytes.Buffer
	tmpl.Execute(&output, params)
	return output.String()
}

func Test_CustomFunctions(t *testing.T) {

	var text = `{{set "name" .personName}}`
	tmpl, err := template.New("").Funcs(TemplatesFunc).Parse(text)
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)

	assert.Equal(t, expand(tmpl, map[string]interface{}{}), "")
	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"personName": "John",
	}), "\"name\":\"John\",")

	tmpl, err = template.New("").
		Funcs(TemplatesFunc).
		Parse(`{{set "age" .person.age}}`)
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"age": 20,
	}), "")

	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"person": map[string]interface{}{
			"age": 20,
		},
	}), "\"age\":20,")

	tmpl, err = template.New("").
		Funcs(TemplatesFunc).
		Parse(`{{set "verbose" .env.VERBOSE}}`)
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"verbose": false,
	}), "")
	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"env": map[string]interface{}{
			"VERBOSE": true,
		},
	}), "\"verbose\":true,")
	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"env": map[string]interface{}{
			"VERBOSE": false,
		},
	}), "\"verbose\":false,")

	tmpl, err = template.New("").
		Funcs(TemplatesFunc).
		Parse(`{{set "verbose" .env.VERBOSE false}}`)
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	assert.Equal(t, expand(tmpl, map[string]interface{}{
		"env": map[string]interface{}{
			"VERBOSE": false,
		},
	}), "\"verbose\":false")

}
