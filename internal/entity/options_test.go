package entity_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/google/go-cmp/cmp"
)

func TestOptions_ToArgs(t *testing.T) {
	tests := []struct {
		name     string
		options  entity.Options
		expected []string
	}{
		{
			name: "boolean true option",
			options: entity.Options{
				"verbose": true,
			},
			expected: []string{"--verbose"},
		},
		{
			name: "boolean false option - should be skipped",
			options: entity.Options{
				"verbose": false,
			},
			expected: []string{},
		},
		{
			name: "string option",
			options: entity.Options{
				"output": "json",
			},
			expected: []string{"--output", "json"},
		},
		{
			name: "int option",
			options: entity.Options{
				"keep-daily": 7,
			},
			expected: []string{"--keep-daily", "7"},
		},
		{
			name: "int64 option",
			options: entity.Options{
				"limit": int64(1000),
			},
			expected: []string{"--limit", "1000"},
		},
		{
			name: "float64 option",
			options: entity.Options{
				"ratio": 0.5,
			},
			expected: []string{"--ratio", "0.5"},
		},
		{
			name: "array option",
			options: entity.Options{
				"exclude": []any{"*.tmp", "node_modules"},
			},
			expected: []string{"--exclude", "*.tmp", "--exclude", "node_modules"},
		},
		{
			name: "multiple options",
			options: entity.Options{
				"verbose":    true,
				"output":     "json",
				"keep-daily": 7,
			},
			expected: []string{"--verbose", "--output", "json", "--keep-daily", "7"},
		},
		{
			name:     "nil options",
			options:  nil,
			expected: []string{},
		},
		{
			name:     "empty options",
			options:  entity.Options{},
			expected: []string{},
		},

		{
			name: "single char short option",
			options: entity.Options{
				"-v": true,
			},
			expected: []string{"-v"},
		},
		{
			name: "short option with value",
			options: entity.Options{
				"-n": "5",
			},
			expected: []string{"-n", "5"},
		},
		{
			name: "multi-char short option stays as is",
			options: entity.Options{
				"-abc": true,
			},
			expected: []string{"-abc"},
		},
		{
			name: "short option false - should be skipped",
			options: entity.Options{
				"-v": false,
			},
			expected: []string{},
		},
		{
			name: "mixed long and short options",
			options: entity.Options{
				"verbose": true,
				"-x":      true,
			},
			expected: []string{"--verbose", "-x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.ToArgs()

			// Sort both slices for comparison since map iteration order is random
			sort.Strings(result)
			sort.Strings(tt.expected)

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("ToArgs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOptions_Merge(t *testing.T) {
	tests := []struct {
		name     string
		base     entity.Options
		other    entity.Options
		expected entity.Options
	}{
		{
			name: "merge non-overlapping options",
			base: entity.Options{
				"verbose": true,
				"output":  "json",
			},
			other: entity.Options{
				"--keep-daily": 7,
			},
			expected: entity.Options{
				"--verbose":    true,
				"--output":     "json",
				"--keep-daily": 7,
			},
		},
		{
			name: "other overrides base",
			base: entity.Options{
				"verbose":  false,
				"--output": "json",
			},
			other: entity.Options{
				"--verbose": true,
				"format":    "text",
			},
			expected: entity.Options{
				"--verbose": true,
				"--output":  "json",
				"--format":  "text",
			},
		},
		{
			name: "merge with empty other",
			base: entity.Options{
				"verbose": true,
			},
			other: entity.Options{},
			expected: entity.Options{
				"--verbose": true,
			},
		},
		{
			name: "merge with empty base",
			base: entity.Options{},
			other: entity.Options{
				"verbose": true,
			},
			expected: entity.Options{
				"--verbose": true,
			},
		},
		{
			name:     "merge two empty options",
			base:     entity.Options{},
			other:    entity.Options{},
			expected: entity.Options{},
		},
		{
			name: "merge with short and long (with -- and without) options",
			base: entity.Options{
				"verbose":  true,
				"--output": "text",
				"-f":       true,
			},
			other: entity.Options{
				"--verbose": false,
				"output":    "json",
				"-f":        false,
			},
			expected: entity.Options{
				"--verbose": false,
				"--output":  "json",
				"-f":        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Merge(tt.other)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Merge() = %v, want %v", result, tt.expected)
			}
		})
	}
}
