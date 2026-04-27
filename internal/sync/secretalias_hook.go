package sync

import (
	"fmt"
	"io"
	"os"
)

// AliasStage returns a pipeline Stage that applies key aliasing to secrets.
func AliasStage(cfg AliasConfig) Stage {
	return Stage{
		Name: "alias",
		Run: func(secrets map[string]string) (map[string]string, error) {
			return ApplyAliases(secrets, cfg)
		},
	}
}

// AliasAndReport applies aliases and writes a summary to w.
// If w is nil, output goes to os.Stdout.
func AliasAndReport(secrets map[string]string, cfg AliasConfig, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		return nil, fmt.Errorf("aliasandreport: secrets map is nil")
	}

	before := len(secrets)
	result, err := ApplyAliases(secrets, cfg)
	if err != nil {
		return nil, err
	}

	applied := 0
	for src, dst := range cfg.Aliases {
		if _, ok := secrets[src]; ok && dst != src {
			applied++
		}
	}

	fmt.Fprintf(w, "alias: %d keys in, %d aliases applied, %d keys out\n", before, applied, len(result))
	return result, nil
}
