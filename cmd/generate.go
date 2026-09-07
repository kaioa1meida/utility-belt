package cmd

import (
	"github.com/spf13/cobra"
)

// NewGenerateCmd creates the parent generate command.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate passwords, UUIDs, and other useful entities",
		Long: `Generate provides commands to generate secure passwords, standard UUIDs (v4/v7),
and other utility data.

Available Commands:
  password    Generate secure random passwords
  uuid        Generate canonical UUIDs (v4 or v7)
`,
	}

	cmd.AddCommand(NewGeneratePasswordCmd())
	cmd.AddCommand(NewGenerateUUIDCmd())

	return cmd
}

