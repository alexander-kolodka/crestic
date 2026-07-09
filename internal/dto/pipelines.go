package dto

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Pipelines []Pipeline

func (ps *Pipelines) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return errors.New("pipelines must be a sequence")
	}

	var out []Pipeline
	for _, elem := range value.Content {
		var p Pipeline
		err := elem.Decode(&p)
		if err != nil {
			return fmt.Errorf("unmarshal pipeline: %w", err)
		}
		out = append(out, p)
	}
	*ps = out
	return nil
}
