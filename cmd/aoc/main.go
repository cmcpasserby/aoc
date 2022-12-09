package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cmcpasserby/aoc"
	"github.com/cmcpasserby/scli"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	longUsage = `aoc is a command line tool for downloading Advent of Code puzzle inputs
any flag not provided will fallback to using values from .aocConfig file or fallback to using current date for -year and -day`
)

var (
	gFlagYear          int
	gFlagDay           int
	gFlagOutput        string
	gFlagSessionCookie string
)

func main() {
	rootFlags := flag.NewFlagSet("aoc", flag.ExitOnError)
	rootFlags.IntVar(&gFlagYear, "year", 0, "year to download from, not not defined will fallback to a year set in the config if present then current year")
	rootFlags.IntVar(&gFlagDay, "day", 0, "day to download from, not not defined will fallback to a day set in the config if present then current day")
	rootFlags.StringVar(&gFlagOutput, "output", "", "defines output path for the downloaded puzzle input, accepts {year} and {day} as placeholders, if no output is provided results go to stdout")
	rootFlags.StringVar(&gFlagSessionCookie, "sessionCookie", "", "sets the session cookie, if not defined session cookie is read from .aocConfig")

	cmd := &scli.Command{
		Usage:         "aoc",
		ShortHelp:     "Tool for downloading Advent of Code Puzzle input data",
		LongHelp:      longUsage,
		FlagSet:       rootFlags,
		ArgsValidator: scli.NoArgs(),
		Exec: func(ctx context.Context, args []string) error {
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

	if err := cmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func write(year, day int, data []byte) error {
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

	if data[len(data)-1] == 0xA {
		data = data[:len(data)-1]
	}

	if _, err := outputWriter.Write(data); err != nil {
		return err
	}

	return nil
}
