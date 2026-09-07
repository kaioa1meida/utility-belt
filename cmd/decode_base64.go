package cmd

import (
	"fmt"

	"github.com/kaioa1meida/utility-belt/internal/codecs"
	"github.com/spf13/cobra"
)

// NewDecodeBase64Cmd creates the decode base64 subcommand.
func NewDecodeBase64Cmd() *cobra.Command {
	var stringFlag string

	cmd := &cobra.Command{
		Use:   "base64",
		Short: "Decode standard Base64 encoded data",
		Long: `Decode standard RFC 4648 Base64 encoded data back to its original representation.

Input can be provided via standard input (pipe) or through the --string flag.

Examples:
  # Decode text via stdin pipe
  echo -n "aGVsbG8=" | utility-belt decode base64

  # Decode text via --string flag
  utility-belt decode base64 --string "aGVsbG8="
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := resolveInput(cmd, stringFlag, true)
			if err != nil {
				return err
			}

			decoded, err := codecs.DecodeBase64(string(input))
			if err != nil {
				return fmt.Errorf("failed to decode base64: %w", err)
			}

			if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(decoded)); err != nil {
				return fmt.Errorf("failed writing decoded output: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&stringFlag, "string", "", "Base64 encoded string to decode (alternative to stdin)")

	return cmd
}
