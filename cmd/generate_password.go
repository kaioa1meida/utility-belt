package cmd

import (
	"fmt"

	"github.com/kaioa1meida/utility-belt/internal/generators"
	"github.com/spf13/cobra"
)

// NewGeneratePasswordCmd creates the generate password subcommand.
func NewGeneratePasswordCmd() *cobra.Command {
	var opts generators.PasswordOptions

	cmd := &cobra.Command{
		Use:   "password",
		Short: "Generate secure random passwords",
		Long: `Generate one or more cryptographically secure random passwords of configurable
length and character sets using crypto/rand.

Supported Character Sets:
  - Uppercase letters: A-Z
  - Lowercase letters: a-z
  - Numbers: 0-9
  - Symbols: !@#$%^&*()-_=+[]{}|;:,.<>?/~

Examples:
  # Generate a single 16-character password (default)
  utility-belt generate password

  # Generate a 32-character password without symbols
  utility-belt generate password --length 32 --no-symbols

  # Generate 5 passwords, one per line
  utility-belt generate password --count 5

  # Generate a numeric PIN of 6 digits
  utility-belt generate password --length 6 --no-symbols --no-uppercase --no-lowercase
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			passwords, err := generators.GeneratePasswords(opts)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			for _, pwd := range passwords {
				if _, err := fmt.Fprintln(out, pwd); err != nil {
					return fmt.Errorf("failed writing password to stdout: %w", err)
				}
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.Length, "length", 16, "Exact length of each password (must be > 0)")
	flags.IntVar(&opts.Count, "count", 1, "Number of passwords to generate (one per line, must be > 0)")
	flags.BoolVar(&opts.NoSymbols, "no-symbols", false, "Exclude symbols (!@#$%^&*()-_=+[]{}|;:,.<>?/~) from the password")
	flags.BoolVar(&opts.NoNumbers, "no-numbers", false, "Exclude numbers (0-9) from the password")
	flags.BoolVar(&opts.NoUppercase, "no-uppercase", false, "Exclude uppercase letters (A-Z) from the password")
	flags.BoolVar(&opts.NoLowercase, "no-lowercase", false, "Exclude lowercase letters (a-z) from the password")

	return cmd
}

