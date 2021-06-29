package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefinitionTOJSON(t *testing.T){
	expected:=`"synthesis_definitions":[{"rules":[{"name":"rule1"},{"name":"rule2"}],"type":"type1"},{"name":"rule3","type":"type2"}]`
	d:=&Definition{}
	d.AddEntry("type1",map[string]interface{}{
		"rules": []map[string]interface{}{
			{"name": "rule1"},
			{"name": "rule2"},
		},
	})
	d.AddEntry("type2",map[string]interface{}{
		"name": "rule3",
	})
	result,err:=d.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, expected,result)

}