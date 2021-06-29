package synthesis

import (
	"bytes"
	"io"

	"github.com/newrelic/nri-config-generator/synthesis/internal"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ProcessSynthesisDefinitions(in []byte) (string, error) {
	definition := &internal.Definition{}
	r := bytes.NewReader(in)
	if err := processYamlWithMultipleDocuments(r, definition); err != nil {
		return "", err
	}
	return definition.ToJSON()
}

func processYamlWithMultipleDocuments(r *bytes.Reader, def *internal.Definition) error {
	decoder := yaml.NewDecoder(r)
	for {
		var document map[string]interface{}
		if err := decoder.Decode(&document); err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if err := addEntryToDefinition(document, def); err != nil {
			return err
		}
	}
	return nil
}

func addEntryToDefinition(doc map[string]interface{}, def *internal.Definition) error {
	entityType, ok := doc["type"]
	if !ok {
		return errors.New("missing required field 'type' in definition")

	}
	sType := entityType.(string)
	details, ok := doc["synthesis"]
	if !ok {
		def.AddEntry(sType, nil)
		return nil
	}
	detailsMap, ok := details.(map[string]interface{})
	if !ok {
		def.AddEntry(sType, nil)
		return nil
	}
	def.AddEntry(sType, detailsMap)
	return nil

}
