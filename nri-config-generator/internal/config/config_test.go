package config

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const text = `
parent:
  child:
    - element1: a
    - element2: b
powerdns_url: http://powerdns:8080/api/v1/
exporter_port: 9120
api_key: 11111-222222-33333-44444
transformations:
 - ignore_metrics:
     - prefixes:
       - kube_daemonset_
 - ignore_metrics:
     - prefixes:
         - kube_daemonset_

`

func TestGeneratorConfigPortAlreadyInUse(t *testing.T) {
	dir := t.TempDir()
	p := path.Join(dir, "config.yml")
	err := os.WriteFile(p, []byte(text), fs.ModePerm)
	assert.NoError(t, err)
	c, _, err := getConfig(&ArgumentList{ConfigPath: p})
	assert.NoError(t, err)
	assert.Len(t, c["transformations"], 2)
	assert.Equal(t, c["exporter_port"], 9120)
}
