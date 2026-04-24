package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

func init() {
	var (
		tmplText  string
		outputKey string
		envFile   string
	)

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Render secrets through a Go template and print the result",
		Long: `Reads secrets from a .env file and renders them through a Go template.
The template receives the secrets as a map[string]string.

Example:
  vaultpull template --text 'export DB={{ index . "DB_URL" }}' --file .env`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if tmplText == "" {
				return fmt.Errorf("--text is required")
			}
			if envFile == "" {
				return fmt.Errorf("--file is required")
			}

			reader, err := sync.NewReader(envFile)
			if err != nil {
				return fmt.Errorf("open env file: %w", err)
			}
			secrets, err := reader.Read()
			if err != nil {
				return fmt.Errorf("read env file: %w", err)
			}

			r, err := sync.NewTemplateRenderer(tmplText)
			if err != nil {
				return fmt.Errorf("parse template: %w", err)
			}

			if outputKey != "" {
				result, err := r.RenderToMap(secrets, outputKey)
				if err != nil {
					return err
				}
				fmt.Fprintf(os.Stdout, "%s=%s\n", outputKey, result[outputKey])
				return nil
			}

			out, err := r.Render(secrets)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&tmplText, "text", "", "Go template text to render")
	cmd.Flags().StringVar(&outputKey, "output-key", "", "If set, store rendered value under this key and print as KEY=VALUE")
	cmd.Flags().StringVar(&envFile, "file", ".env", "Path to the .env file to read secrets from")

	rootCmd.AddCommand(cmd)
}
