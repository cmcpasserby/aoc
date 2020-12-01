package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

const version = "1.0.0"

func Execute() {
	cmd := &cobra.Command{
		Use:     "aoc [sub]",
		Version: version,
		Short:   "Tool for pulling Advent of Code inputs and submitting answerers",
	}

	cmd.AddCommand(
		createGetCmd(),
		createSubmitCmd(),
	)

	// falls back to get, of no subcommand is provided
	targetCmd, _, err := cmd.Find(os.Args[1:])
	if err != nil || targetCmd == nil {
		args := append([]string{"get"}, os.Args[1:]...)
		cmd.SetArgs(args)
	}

	_ = cmd.Execute()
}
