package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Fetch a random sample of secrets from Vault",
	RunE:  runSample,
}

func runSample(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	maxSamples, _ := cmd.Flags().GetInt("max")
	seed, _ := cmd.Flags().GetInt64("seed")
	onlyKeys, _ := cmd.Flags().GetStringSlice("only-keys")

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("sample: vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("sample: fetch secrets: %w", err)
	}

	sampleCfg := sync.DefaultSampleConfig()
	if maxSamples > 0 {
		sampleCfg.MaxSamples = maxSamples
	}
	if seed != 0 {
		sampleCfg.Seed = seed
	}
	sampleCfg.OnlyKeys = onlyKeys

	sampled, err := sync.SampleSecrets(secrets, sampleCfg)
	if err != nil {
		return fmt.Errorf("sample: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", sync.SampleSummary(len(secrets), len(sampled)))
	for k, v := range sampled {
		fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
	}
	return nil
}

func init() {
	sampleCmd.Flags().Int("max", 10, "maximum number of secrets to sample")
	sampleCmd.Flags().Int64("seed", 42, "random seed for deterministic sampling")
	sampleCmd.Flags().StringSlice("only-keys", nil, "restrict sampling to keys with these prefixes")
	rootCmd.AddCommand(sampleCmd)
}
