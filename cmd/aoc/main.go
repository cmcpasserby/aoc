package main

import (
	"fmt"
	"github.com/cmcpasserby/aoc"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	version   = "1.0.1"
	longUsage = `aoc is a command line tool for downloading Advent of Code puzzle inputs
any flag not provided will fallback to using values from .aocConfig file or fallback to using current date for --year and --day`
)

func main() {
	var (
		lFlagYear          int
		lFlagDay           int
		lFlagOutput        string
		lFlagSessionCookie string
	)

	cmd := &cobra.Command{
		Use:           "aoc",
		Version:       version,
		Short:         "Tool for downloading Advent of Code puzzle input data",
		Long:          longUsage,
		Args:          cobra.NoArgs,
		SilenceUsage:  false,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := getConfig()
			if err != nil {
				return err
			}

			sessionCookie := config.SessionCookie
			if lFlagSessionCookie != "" {
				sessionCookie = lFlagSessionCookie
			}

			now := time.Now()
			year := now.Year()
			if lFlagYear != 0 {
				year = lFlagYear
			} else if config.Year != 0 {
				year = config.Year
			}

			day := now.Day()
			if lFlagDay != 0 {
				day = lFlagDay
			} else if config.Day != 0 {
				day = config.Day
			}

			client := aoc.New(sessionCookie, http.DefaultClient)

			result, err := client.DownloadInput(year, day)
			if err != nil {
				return err
			}


			var outputWriter io.Writer

			if lFlagOutput != "" {
				output := strings.ReplaceAll(lFlagOutput, "{{year}}", fmt.Sprintf("%04d", year))
				output = strings.ReplaceAll(output, "{{day}}", fmt.Sprintf("%02d", day))
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
				outputWriter = f
			} else {
				outputWriter = os.Stdout
			}

			if _, err = io.Copy(outputWriter, result); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&lFlagYear, "year", "y", 0, "year to download from, if not defined will fallback to year set in config if present then the current year")
	cmd.Flags().IntVarP(&lFlagDay, "day", "d", 0, "day to download from, if not defined will fallback to day set in config if present then current year")
	cmd.Flags().StringVarP(&lFlagOutput, "output", "o", "", "defines output path for downloaded puzzle input, accepts {{year}} and {{day}} as place holders for year and day, if not output is provided result goes to stdout")
	cmd.Flags().StringVar(&lFlagSessionCookie, "sessionCookie", "", "sets the session cookie, if not defined session cookie is read from .aocConfig")

	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}
