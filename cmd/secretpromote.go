package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Promote secrets from one prefix to another",
	Long: `Copies secrets matching --from-prefix into --to-prefix.
By default existing destination keys are not overwritten.`,
	RunE: runPromote,
}

func runPromote(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fromPrefix, _ := cmd.Flags().GetString("from-prefix")
	toPrefix, _ := cmd.Flags().GetString("to-prefix")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	vc, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vc.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	pcfg := sync.PromoteConfig{
		FromPrefix: fromPrefix,
		ToPrefix:   toPrefix,
		Overwrite:  overwrite,
		DryRun:     dryRun,
	}

	_, err = sync.PromoteAndReport(secrets, pcfg, os.Stdout)
	return err
}

func init() {
	promoteCmd.Flags().String("from-prefix", "", "Source prefix to match (e.g. dev_)")
	promoteCmd.Flags().String("to-prefix", "", "Destination prefix to write (e.g. prod_)")
	_ = promoteCmd.MarkFlagRequired("to-prefix")
	promoteCmd.Flags().Bool("overwrite", false, "Overwrite existing destination keys")
	promoteCmd.Flags().Bool("dry-run", false, "Preview promotions without writing")
	rootCmd.AddCommand(promoteCmd)
}
