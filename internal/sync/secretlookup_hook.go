package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// LookupAndReport performs a lookup and writes a formatted report to w.
// If w is nil, os.Stdout is used. Returns the matching results.
func LookupAndReport(secrets map[string]string, queries []string, cfg LookupConfig, w io.Writer) ([]LookupResult, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		return nil, fmt.Errorf("lookup: secrets map is nil")
	}

	results, err := LookupSecrets(secrets, queries, cfg)
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	fmt.Fprintln(w, LookupSummary(results))
	for _, r := range results {
		fmt.Fprintf(w, "  %s = %s\n", r.Key, r.Value)
	}
	return results, nil
}

// LookupStage returns a pipeline stage that injects lookup metadata.
// For each query, it adds a key "__lookup_<query>" with the matched value (or "" if not found).
func LookupStage(queries []string, cfg LookupConfig) PipelineStage {
	return PipelineStage{
		Name: "lookup",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			if secrets == nil {
				return nil, fmt.Errorf("lookup stage: nil secrets")
			}
			out := make(map[string]string, len(secrets))
			for k, v := range secrets {
				out[k] = v
			}
			results, err := LookupSecrets(secrets, queries, cfg)
			if err != nil {
				return nil, err
			}
			hit := make(map[string]string, len(results))
			for _, r := range results {
				hit[r.Key] = r.Value
			}
			for _, q := range queries {
				annotationKey := fmt.Sprintf("__lookup_%s", q)
				if v, ok := hit[q]; ok {
					out[annotationKey] = v
				} else {
					out[annotationKey] = ""
				}
			}
			return out, nil
		},
	}
}
