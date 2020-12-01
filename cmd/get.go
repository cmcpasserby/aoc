package cmd

import (
	"fmt"
	"github.com/cmcpasserby/aoc/aoc"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"time"
)

func createGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get [year day]",
		Short: "gets inputs for a given year and day",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) > 2 {
				return fmt.Errorf("accepts 0 or 2 args, received %d", len(args))
			}
			
			if len(args) == 0 {
				return nil
			}

			now := time.Now()

			year, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if year < 2015 || year > now.Year() {
				return fmt.Errorf("please choose a valid year (range 2015 - %d), received %d", now.Year(), year)
			}

			day, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			if day < 1 || day > 31 {
				return fmt.Errorf("please choose a valid day (range 1 - 31), received %d", day)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			year, day, err := parseInputsArgs(args)
			if err != nil {
				return err
			}

			config, err := loadConfig()
			if err != nil {
				return err
			}

			client, err := aoc.NewClient(config.SessionId)
			if err != nil {
				return err
			}

			inputs, err := client.GetInputs(year, day)
			if err != nil {
				return err
			}

			_, err = io.Copy(os.Stdout, inputs.Input)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func parseInputsArgs(args []string) (year, day int, err error) {
	if len(args) == 0 {
		now := time.Now()
		return now.Year(), now.Day(), nil
	}

	year, err = strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, err
	}

	day, err = strconv.Atoi(args[1])
	if err != nil {
		return 0, 0, err
	}

	return year, day, nil
}

