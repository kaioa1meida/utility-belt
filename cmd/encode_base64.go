package cmd

import (
	"fmt"

	"github.com/kaioa1meida/utility-belt/internal/codecs"
	"github.com/spf13/cobra"
)

// NewEncodeBase64Cmd creates the encode base64 subcommand.
func NewEncodeBase64Cmd() *cobra.Command {
	var stringFlag string

	cmd := &cobra.Command{
		Use:   "base64",
		Short: "Encode data into standard Base64 format",
		Long: `Encode data into standard RFC 4648 Base64 format.

Input can be provided via standard input (pipe) or through the --string flag.

Examples:
  # Encode text via stdin pipe
  echo -n "hello" | utility-belt encode base64

  # Encode text via --string flag
  utility-belt encode base64 --string "hello"

  # Copy encoded output to clipboard (macOS)
  utility-belt encode base64 --string "secret" | pbcopy
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := resolveInput(cmd, stringFlag, false)
			if err != nil {
				return err
			}

			encoded := codecs.EncodeBase64(input)
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), encoded); err != nil {
				return fmt.Errorf("failed writing base64 output: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&stringFlag, "string", "", "String to encode into Base64 (alternative to stdin)")

	return cmd
}
