package main

import (
	"fmt"
	"github.com/cmcpasserby/aoc"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	version   = "1.0.1"
	longUsage = `aoc is a command line tool for downloading Advent of Code puzzle inputs
any flag not provided will fallback to using values from .aocConfig file or fallback to using current date for --year and --day`
)

var (
	gFlagYear          int
	gFlagDay           int
	gFlagOutput        string
	gFlagSessionCookie string
)

func main() {
	cmd := &cobra.Command{
		Use:           "aoc",
		Version:       version,
		Short:         "Tool for downloading Advent of Code puzzle input data",
		Long:          longUsage,
	}

	cmd.AddCommand(
		createInputsCmd(),
		createQuestionCmd(),
	)

	cmd.PersistentFlags().IntVarP(&gFlagYear, "year", "y", 0, "year to download from, if not defined will fallback to year set in config if present then the current year")
	cmd.PersistentFlags().IntVarP(&gFlagDay, "day", "d", 0, "day to download from, if not defined will fallback to day set in config if present then current year")
	cmd.PersistentFlags().StringVarP(&gFlagOutput, "output", "o", "", "defines output path for downloaded puzzle input, accepts {{year}} and {{day}} as place holders for year and day, if not output is provided result goes to stdout")
	cmd.PersistentFlags().StringVar(&gFlagSessionCookie, "sessionCookie", "", "sets the session cookie, if not defined session cookie is read from .aocConfig")

	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func write(year, day int, reader io.Reader) error {
	var outputWriter io.Writer

	if gFlagOutput != "" {
		output := strings.ReplaceAll(gFlagOutput, "{year}", fmt.Sprintf("%04d", year))
		output = strings.ReplaceAll(output, "{day}", fmt.Sprintf("%02d", day))
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		outputWriter = f
	} else {
		outputWriter = os.Stdout
	}

	if _, err := io.Copy(outputWriter, reader); err != nil {
		return err
	}

	return nil
}

func createInputsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "inputs",
		Short: "Gets puzzle inputs for a day",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := getCombinedConfig()
			if err != nil {
				return err
			}
			client := aoc.New(config.SessionCookie, http.DefaultClient)
			inputs, err := client.GetInput(config.Year, config.Day)
			if err != nil {
				return err
			}
			return write(config.Year, config.Day, inputs)
		},
	}
}

func createQuestionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "question",
		Short: "Gets question for day",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := getCombinedConfig()
			if err != nil {
				return err
			}

			client := aoc.New(config.SessionCookie, http.DefaultClient)
			question, err := client.GetQuestion(config.Year, config.Day)
			if err != nil {
				return err
			}
			return write(config.Year, config.Day, question)
		},
	}
}
