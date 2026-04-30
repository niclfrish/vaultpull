package sync

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents the target type for casting a secret value.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastBool   CastType = "bool"
)

// CastRule defines a key pattern and the target type to cast its value to.
type CastRule struct {
	Key      string
	CastTo   CastType
}

// DefaultCastConfig returns a default set of cast rules (empty).
func DefaultCastConfig() []CastRule {
	return []CastRule{}
}

// CastSecrets applies type-casting rules to secret values, normalising them
// into canonical string representations (e.g. "true" for booleans, "42" for
// integers). Keys not matched by any rule are returned unchanged.
func CastSecrets(secrets map[string]string, rules []CastRule) (map[string]string, error) {
	if secrets == nil {
		return nil, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for _, rule := range rules {
		v, ok := out[rule.Key]
		if !ok {
			continue
		}

		casted, err := castValue(v, rule.CastTo)
		if err != nil {
			return nil, fmt.Errorf("cast %q to %s: %w", rule.Key, rule.CastTo, err)
		}
		out[rule.Key] = casted
	}

	return out, nil
}

func castValue(v string, t CastType) (string, error) {
	v = strings.TrimSpace(v)
	switch t {
	case CastString:
		return v, nil
	case CastInt:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid integer %q", v)
		}
		return strconv.FormatInt(n, 10), nil
	case CastFloat:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "", fmt.Errorf("invalid float %q", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CastBool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", fmt.Errorf("invalid boolean %q", v)
		}
		return strconv.FormatBool(b), nil
	default:
		return "", fmt.Errorf("unknown cast type %q", t)
	}
}

// CastSummary returns a human-readable summary of applied cast rules.
func CastSummary(rules []CastRule) string {
	if len(rules) == 0 {
		return "no cast rules applied"
	}
	return fmt.Sprintf("%d cast rule(s) applied", len(rules))
}
