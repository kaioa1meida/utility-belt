package cmd

import (
	"github.com/spf13/cobra"
)

// NewDecodeCmd creates the parent decode command.
func NewDecodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode",
		Short: "Decode data from various formats",
		Long: `Decode provides commands to decode data from formats such as Base64 and JWT.

Available Commands:
  base64    Decode standard Base64 encoded data
  jwt       Decode JSON Web Token (Header and Payload)
`,
	}

	cmd.AddCommand(NewDecodeBase64Cmd())
	cmd.AddCommand(NewDecodeJWTCmd())

	return cmd
}
