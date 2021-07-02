package synthesis

import (
	"bytes"
	"errors"
	"io"

	"github.com/newrelic/nri-config-generator/synthesis/internal"
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
	sType, details, err := getAttributesFromSynthesis(doc)
	if err != nil {
		return err
	}
	def.AddEntry(sType, details)
	return nil
}

func getAttributesFromSynthesis(doc map[string]interface{}) (string, map[string]interface{}, error) {
	sType, ok := doc["type"]
	if !ok {
		return "", nil, errors.New("missing required field 'type' in definition")
	}
	sTypeStr, ok := sType.(string)
	if !ok {
		return "", nil, errors.New("invalid format for field 'type', It must be a string")
	}
	synthesis, ok := doc["synthesis"]
	if !ok {
		return "", nil, errors.New("missing required field 'synthesis' in definition")
	}
	synthesisMap, ok := synthesis.(map[string]interface{})
	if !ok {
		return "", nil, errors.New("invalid format for field 'synthesis', It must be an object")
	}
	return sTypeStr, synthesisMap, nil
}
