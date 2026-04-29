package sync

import (
	"fmt"
	"strings"
)

// MergeStrategy defines how conflicting keys are resolved when merging secret maps.
type MergeStrategy string

const (
	MergeStrategyFirst  MergeStrategy = "first"  // keep the first value seen
	MergeStrategyLast   MergeStrategy = "last"   // keep the last value seen
	MergeStrategyError  MergeStrategy = "error"  // return an error on conflict
	MergeStrategyPrefix MergeStrategy = "prefix" // prefix conflicting keys with source index
)

// DefaultMergeStrategyConfig returns a config using the "last" strategy.
func DefaultMergeStrategyConfig() MergeStrategyConfig {
	return MergeStrategyConfig{
		Strategy: MergeStrategyLast,
	}
}

// MergeStrategyConfig controls how MergeWithStrategy behaves.
type MergeStrategyConfig struct {
	Strategy MergeStrategy
}

// MergeWithStrategy merges multiple secret maps according to the configured strategy.
// Sources are applied in order; later sources may override earlier ones depending on strategy.
func MergeWithStrategy(cfg MergeStrategyConfig, sources ...map[string]string) (map[string]string, error) {
	if len(sources) == 0 {
		return map[string]string{}, nil
	}

	switch cfg.Strategy {
	case MergeStrategyFirst:
		return mergeFirst(sources)
	case MergeStrategyLast:
		return mergeLast(sources)
	case MergeStrategyError:
		return mergeError(sources)
	case MergeStrategyPrefix:
		return mergePrefix(sources)
	default:
		return nil, fmt.Errorf("unknown merge strategy: %q", cfg.Strategy)
	}
}

func mergeFirst(sources []map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for _, src := range sources {
		for k, v := range src {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
	}
	return out, nil
}

func mergeLast(sources []map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for _, src := range sources {
		for k, v := range src {
			out[k] = v
		}
	}
	return out, nil
}

func mergeError(sources []map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for i, src := range sources {
		for k, v := range src {
			if existing, exists := out[k]; exists && existing != v {
				return nil, fmt.Errorf("conflict on key %q at source index %d", k, i)
			}
			out[k] = v
		}
	}
	return out, nil
}

func mergePrefix(sources []map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for i, src := range sources {
		prefix := fmt.Sprintf("SRC%d_", i)
		for k, v := range src {
			key := strings.ToUpper(k)
			if _, exists := out[key]; exists {
				out[prefix+key] = v
			} else {
				out[key] = v
			}
		}
	}
	return out, nil
}
