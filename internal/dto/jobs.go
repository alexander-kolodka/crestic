package dto

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type (
	Job  any
	Jobs []Job
)

func (js *Jobs) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return errors.New("jobs must be a sequence")
	}

	var out []Job
	for _, elem := range value.Content {
		var w jobWrapper
		err := elem.Decode(&w)
		if err != nil {
			return fmt.Errorf("unmarshal job: %w", err)
		}
		out = append(out, w.Job)
	}
	*js = out
	return nil
}

type jobWrapper struct {
	Job Job
}

func (w *jobWrapper) UnmarshalYAML(value *yaml.Node) error {
	var t struct {
		Type string `yaml:"type"`
	}
	err := value.Decode(&t)
	if err != nil {
		return fmt.Errorf("job type: %w", err)
	}

	switch t.Type {
	case "backup":
		var b BackupJob
		dErr := value.Decode(&b)
		if dErr != nil {
			return fmt.Errorf("backup: %w", dErr)
		}
		w.Job = b
	case "copy":
		var c CopyJob
		dErr := value.Decode(&c)
		if dErr != nil {
			return fmt.Errorf("copy: %w", dErr)
		}
		w.Job = c
	default:
		return fmt.Errorf("unknown job type: %q", t.Type)
	}
	return nil
}
