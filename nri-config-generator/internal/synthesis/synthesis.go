package synthesis

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	attType = "type"
)

type Definition []synthesis

type synthesis map[string]interface{}

func (d *Definition) addEntry(sType string, details map[string]interface{}) {
	s := synthesis{
		attType: sType,
	}

	// copying details into synthesis struct.
	for k, v := range details {
		s[k] = v
	}
	*d = append(*d, s)
}

func (d *Definition) ToJSON() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ProcessSynthesisDefinitions(in []byte) (string, error) {
	definition := &Definition{}
	r := bytes.NewReader(in)
	if err := processYamlWithMultipleDocuments(r, definition); err != nil {
		return "", err
	}

	return definition.ToJSON()
}

func processYamlWithMultipleDocuments(r *bytes.Reader, def *Definition) error {
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

func addEntryToDefinition(doc map[string]interface{}, def *Definition) error {
	sType, details, err := getAttributesFromSynthesis(doc)
	if err != nil {
		return err
	}
	def.addEntry(sType, details)

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
