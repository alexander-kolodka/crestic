package entity

import (
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// Options is a map of command line options.
type Options map[string]any

// ToArgs converts options map to command line arguments
// Supports both long options (--key) and short options (-k)
// If key starts with "-", it's treated as a short option.
func (opts Options) ToArgs() []string {
	return lo.Flatten(lo.MapToSlice(opts, func(flag string, value any) []string {
		flag = formatFlag(flag)

		b, ok := value.(bool)
		if ok {
			return toBoolArgs(flag, b)
		}

		arr, ok := value.([]any)
		if ok {
			return toArrayArgs(flag, arr)
		}

		return toScalarArgs(flag, value)
	}))
}

// Merge merges current options with another options map (other overrides current).
func (opts Options) Merge(other Options) Options {
	result := make(Options)
	for f, v := range opts {
		result[formatFlag(f)] = v
	}
	for f, v := range other {
		result[formatFlag(f)] = v
	}
	return result
}

// valueToString converts any value to string.
func valueToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return ""
	}
}

func formatFlag(flag string) string {
	if hasFlagHyphen(flag) {
		return flag
	}

	return "--" + flag
}

func hasFlagHyphen(key string) bool {
	return strings.HasPrefix(key, "-")
}

func toArrayArgs(flag string, values []any) []string {
	return lo.Flatten(lo.Map(values, func(v any, _ int) []string {
		return toScalarArgs(flag, v)
	}))
}

func toBoolArgs(flag string, value bool) []string {
	if !value {
		return nil
	}

	return []string{flag}
}

func toScalarArgs(flag string, value any) []string {
	v := valueToString(value)
	if v == "" {
		return nil
	}

	return []string{flag, v}
}
