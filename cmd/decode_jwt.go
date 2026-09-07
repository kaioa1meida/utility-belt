package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kaioa1meida/utility-belt/internal/codecs"
	"github.com/spf13/cobra"
)

// NewDecodeJWTCmd creates the decode jwt subcommand.
func NewDecodeJWTCmd() *cobra.Command {
	var stringFlag string

	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "Decode JSON Web Token (Header and Payload)",
		Long: `Decode a JSON Web Token (JWS compact format) into structured JSON containing
its Header and Payload.

This command is offline and stateless: it does NOT verify cryptographic signatures.
If the token contains an 'exp' claim that has passed, a warning is printed to stderr.

Examples:
  # Decode JWT via stdin pipe
  echo -n "eyJhbGciOi..." | utility-belt decode jwt

  # Decode JWT via --string flag
  utility-belt decode jwt --string "eyJhbGciOi..."
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := resolveInput(cmd, stringFlag, true)
			if err != nil {
				return err
			}

			result, err := codecs.DecodeJWT(string(input), time.Now())
			if err != nil {
				return fmt.Errorf("failed to decode jwt: %w", err)
			}

			if result.IsExpired {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning: token is expired")
			}

			outJSON, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format JSON output: %w", err)
			}

			if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(outJSON)); err != nil {
				return fmt.Errorf("failed writing jwt output: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&stringFlag, "string", "", "JWT token string to decode (alternative to stdin)")

	return cmd
}
