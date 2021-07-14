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
