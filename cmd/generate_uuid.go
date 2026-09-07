package cmd

import (
	"fmt"

	"github.com/kaioa1meida/utility-belt/internal/generators"
	"github.com/spf13/cobra"
)

// NewGenerateUUIDCmd creates the generate uuid subcommand.
func NewGenerateUUIDCmd() *cobra.Command {
	var opts generators.UUIDOptions

	cmd := &cobra.Command{
		Use:   "uuid",
		Short: "Generate canonical UUIDs (v4 or v7)",
		Long: `Generate canonical 8-4-4-4-12 lowercase UUIDs.

By default, UUID version 4 (random) is generated. If --v7 is specified,
UUID version 7 (timestamp-based, RFC 9562) is generated.

Examples:
  # Generate a single UUID v4 (default)
  utility-belt generate uuid

  # Generate a single UUID v7
  utility-belt generate uuid --v7

  # Generate 3 UUIDs, one per line
  utility-belt generate uuid --count 3

  # Generate 5 UUIDs v7
  utility-belt generate uuid --v7 --count 5
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			uuids, err := generators.GenerateUUIDs(opts)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			for _, id := range uuids {
				if _, err := fmt.Fprintln(out, id); err != nil {
					return fmt.Errorf("failed writing uuid to stdout: %w", err)
				}
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.V7, "v7", false, "Generate UUID version 7 (RFC 9562 timestamp-based) instead of version 4")
	flags.IntVar(&opts.Count, "count", 1, "Number of UUIDs to generate (one per line, must be > 0)")

	return cmd
}
