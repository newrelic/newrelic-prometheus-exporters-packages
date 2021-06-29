package synthesis

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessSynthesisDefinitionsNoSynthesisBlock(t *testing.T){
	in,err:=ioutil.ReadFile("testdata/sample.yml")
	assert.Nil(t, err)
	res,err:=ProcessSynthesisDefinitions(in)
	assert.Nil(t, err)
	assert.Equal(t,  "\"synthesis_definitions\":[{\"type\":\"RAVENDB_DATABASE\"},{\"type\":\"RAVENDB_NODE\"}]", res)
}

func Test_ProcessSynthesisDefinitions(t *testing.T){
	in,err:=ioutil.ReadFile("testdata/sample2.yml")
	assert.Nil(t, err)
	res,err:=ProcessSynthesisDefinitions(in)
	assert.Nil(t, err)
	assert.Equal(t,  "\"synthesis_definitions\":[{\"conditions\":[{\"attribute\":\"metricName\",\"prefix\":\"redis_\"}],\"encodeIdentifierInGUID\":true,\"identifier\":\"targetName\",\"name\":\"targetName\",\"tags\":{\"clusterName\":null,\"targetName\":null},\"type\":\"REDIS\"}]", res)
}