package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var (
	cacheTTL      time.Duration
	cacheDir      string
	cacheInvalidate bool
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the local secrets cache",
}

var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache status for a secret path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		ns, _ := cmd.Flags().GetString("namespace")
		if path == "" {
			return fmt.Errorf("--path is required")
		}
		c, err := sync.NewSecretCache(resolvedCacheDir())
		if err != nil {
			return err
		}
		entry, err := c.Get(path, ns)
		if err != nil {
			return err
		}
		if entry == nil {
			fmt.Println("cache: no entry found")
			return nil
		}
		age := time.Since(entry.FetchedAt).Round(time.Second)
		fmt.Printf("cache: path=%s namespace=%s keys=%d age=%s checksum=%s\n",
			entry.Path, entry.Namespace, len(entry.Secrets), age, entry.Checksum[:12])
		return nil
	},
}

var cacheInvalidateCmd = &cobra.Command{
	Use:   "invalidate",
	Short: "Invalidate cached secrets for a path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		ns, _ := cmd.Flags().GetString("namespace")
		if path == "" {
			return fmt.Errorf("--path is required")
		}
		c, err := sync.NewSecretCache(resolvedCacheDir())
		if err != nil {
			return err
		}
		if err := c.Invalidate(path, ns); err != nil {
			return err
		}
		fmt.Printf("cache: invalidated %s (namespace: %q)\n", path, ns)
		return nil
	},
}

func resolvedCacheDir() string {
	if cacheDir != "" {
		return cacheDir
	}
	home, _ := os.UserCacheDir()
	return filepath.Join(home, "vaultpull", "cache")
}

func init() {
	for _, sub := range []*cobra.Command{cacheStatusCmd, cacheInvalidateCmd} {
		sub.Flags().String("path", "", "Vault secret path")
		sub.Flags().String("namespace", "", "Vault namespace")
	}
	cacheCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "Override cache directory")
	cacheCmd.AddCommand(cacheStatusCmd, cacheInvalidateCmd)
	rootCmd.AddCommand(cacheCmd)
}
