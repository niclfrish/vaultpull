package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

// multienvCmd writes vault secrets to multiple .env targets in one pass.
var multienvCmd = &cobra.Command{
	Use:   "multienv",
	Short: "Write secrets to multiple .env targets simultaneously",
	RunE:  runMultiEnv,
}

func runMultiEnv(cmd *cobra.Command, _ []string) error {
	targetFlag, _ := cmd.Flags().GetStringArray("target")
	if len(targetFlag) == 0 {
		return fmt.Errorf("at least one --target is required (format: name:path[:namespace])")
	}

	targets, err := parseTargets(targetFlag)
	if err != nil {
		return err
	}

	client, err := vault.New(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"))
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secretPath := os.Getenv("VAULT_SECRET_PATH")
	if secretPath == "" {
		return fmt.Errorf("VAULT_SECRET_PATH is required")
	}

	secrets, err := client.GetSecrets(cmd.Context(), secretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	writerFn := func(path string, s map[string]string) error {
		w, err := env.New(path)
		if err != nil {
			return err
		}
		return w.Write(s)
	}

	mw, err := sync.NewMultiEnvWriter(targets, writerFn)
	if err != nil {
		return err
	}

	results := mw.WriteAll(secrets)
	for name, werr := range results {
		if werr != nil {
			fmt.Fprintf(os.Stderr, "target %q failed: %v\n", name, werr)
		} else {
			fmt.Fprintf(os.Stdout, "target %q written\n", name)
		}
	}
	return sync.AnyError(results)
}

func parseTargets(raw []string) ([]sync.EnvTarget, error) {
	var targets []sync.EnvTarget
	for _, r := range raw {
		parts := strings.SplitN(r, ":", 3)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid target %q: expected name:path[:namespace]", r)
		}
		t := sync.EnvTarget{Name: parts[0], Path: parts[1]}
		if len(parts) == 3 {
			t.Namespace = parts[2]
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func init() {
	multienvCmd.Flags().StringArray("target", nil, "Output target in format name:path[:namespace] (repeatable)")
	rootCmd.AddCommand(multienvCmd)
}
