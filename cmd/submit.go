package cmd

import (
	"fmt"
	"github.com/cmcpasserby/aoc/aoc"
	"github.com/spf13/cobra"
	"strconv"
)

func createSubmitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "submit [year day part answer]",
		Short: "submits answer for a day and part",
		Args: cobra.ExactValidArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadConfig()
			if err != nil {
				return err
			}

			client, err := aoc.NewClient(config.SessionId)
			if err != nil {
				return err
			}

			answer, err := parseSubmitArgs(args)
			if err != nil {
				return err
			}

			err, msg := client.Submit(answer)
			if err != nil {
				return err
			}

			fmt.Println(msg)

			return nil
		},
	}

	return cmd
}

func parseSubmitArgs(args []string) (*aoc.Answer, error) {
	year, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	day, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}

	part, err := aoc.ParsePuzzlePart(args[2])
	if err != nil {
		return nil, err
	}

	a := &aoc.Answer{
		Year: year,
		Day: day,
		Part: part,
		Answer: args[3],
	}

	return a, nil
}
