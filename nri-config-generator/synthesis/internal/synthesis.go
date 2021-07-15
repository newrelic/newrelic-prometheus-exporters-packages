package internal

import (
	"encoding/json"
)

const (
	attType = "type"
)

type Definition []synthesis

type synthesis map[string]interface{}

func newSynthesis(sType string) synthesis {
	return synthesis{
		attType: sType,
	}
}

func (s synthesis) withDetails(attributes map[string]interface{}) synthesis {
	for k, v := range attributes {
		s[k] = v
	}
	return s
}

func (d *Definition) AddEntry(sType string, details map[string]interface{}) {
	*d = append(*d, newSynthesis(sType).withDetails(details))
}

func (d *Definition) ToJSON() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
