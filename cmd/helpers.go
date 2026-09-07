package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var ErrNoInput = errors.New("no input provided: provide input via stdin or --string flag")

// resolveInput returns data from the --string flag if present, or reads from stdin.
// When trimText is true, leading and trailing whitespace is stripped (useful for decode operations).
func resolveInput(cmd *cobra.Command, stringFlag string, trimText bool) ([]byte, error) {
	if stringFlag != "" {
		if trimText {
			stringFlag = strings.TrimSpace(stringFlag)
		}
		return []byte(stringFlag), nil
	}

	in := cmd.InOrStdin()
	if in == nil {
		return nil, ErrNoInput
	}

	if file, ok := in.(*os.File); ok && file == os.Stdin {
		stat, err := file.Stat()
		if err != nil || (stat.Mode()&os.ModeCharDevice) != 0 {
			return nil, ErrNoInput
		}
	}

	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrNoInput
	}
	if trimText {
		return []byte(strings.TrimSpace(string(data))), nil
	}
	return data, nil
}
