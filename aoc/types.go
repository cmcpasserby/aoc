package aoc

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type SubmitResponse string

type PuzzlePart int

const (
	PuzzlePartA PuzzlePart = 1
	PuzzlePartB PuzzlePart = 2
)

func ParsePuzzlePart(s string) (PuzzlePart, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	part := PuzzlePart(i)
	if part != PuzzlePartA && part != PuzzlePartB {
		return 0, fmt.Errorf("invalid puzzle part, expected (%d or %d), received %d", PuzzlePartA, PuzzlePartB, part)
	}

	return part, nil
}

type Puzzle struct {
	Year  int
	Day   int
	Input io.Reader
}

type Answer struct {
	Year   int
	Day    int
	Part   PuzzlePart
	Answer string
}

func (a *Answer) Validate() error {
	now := time.Now()
	if a.Year < 2015 || a.Year > now.Year() {
		return fmt.Errorf("answer validate: invalid year provided (range 2015 %d), received: %d", now.Year(), a.Year)
	}

	if a.Day < 1 || a.Day > 31 {
		return fmt.Errorf("answer validate: invalid day provided (range 1 31), received: %d", a.Day)
	}

	if a.Answer == "" {
		return fmt.Errorf("answer validate: no answer provided")
	}

	return nil
}
