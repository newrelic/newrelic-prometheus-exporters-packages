package synthesis

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessSynthesisDefinitionsNoSynthesisBlock(t *testing.T) {
	in, err := ioutil.ReadFile("testdata/sample.yml")
	assert.Nil(t, err)
	res, err := ProcessSynthesisDefinitions(in)
	assert.NotNil(t, err)
	assert.Equal(t, "missing required field 'synthesis' in definition", err.Error())
	assert.Empty(t, res)
}

func Test_ProcessSynthesisDefinitions(t *testing.T) {
	in, err := ioutil.ReadFile("testdata/sample2.yml")
	assert.Nil(t, err)
	res, err := ProcessSynthesisDefinitions(in)
	assert.Nil(t, err)
	assert.Equal(t, "[{\"conditions\":[{\"attribute\":\"metricName\",\"prefix\":\"redis_\"}],\"encodeIdentifierInGUID\":true,\"identifier\":\"targetName\",\"name\":\"targetName\",\"tags\":{\"clusterName\":null,\"targetName\":null},\"type\":\"REDIS\"}]", res)
}

func Test_DefinitionTOJSON(t *testing.T) {
	expected := `[{"rules":[{"name":"rule1"},{"name":"rule2"}],"type":"type1"},{"name":"rule3","type":"type2"}]`
	d := &Definition{}
	d.addEntry("type1", map[string]interface{}{
		"rules": []map[string]interface{}{
			{"name": "rule1"},
			{"name": "rule2"},
		},
	})
	d.addEntry("type2", map[string]interface{}{
		"name": "rule3",
	})
	result, err := d.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

}
