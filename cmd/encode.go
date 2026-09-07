package cmd

import (
	"github.com/spf13/cobra"
)

// NewEncodeCmd creates the parent encode command.
func NewEncodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode data into various formats",
		Long: `Encode provides commands to encode data into standard formats such as Base64.

Available Commands:
  base64    Encode data into standard Base64 format
`,
	}

	cmd.AddCommand(NewEncodeBase64Cmd())

	return cmd
}
