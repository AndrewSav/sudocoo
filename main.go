package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AndrewSav/sudocoo/pkg/format"
	"github.com/AndrewSav/sudocoo/pkg/parser"
	"github.com/AndrewSav/sudocoo/pkg/solver"
)

func main() {

	flags := ParseArgs()

	scanner := parser.CreateInputScanner(flags.InputReader)

	// Statistics block
	totalSolutions := 0
	globalLimit := false
	puzzleCount := 0
	iterations := 0
	start := time.Now()

	for ; ; puzzleCount++ {
		puzzleInput, err := parser.ReadNextPuzzleInput(scanner)
		// if this is the first puzzle and there is no puzzle,
		// then it's a error, otherwise we processed all puzzles
		if errors.Is(err, io.EOF) && puzzleCount != 0 {
			break
		}
		// we exit on these errors because the format is realy loose
		// and it is unlikely we can recover once something went wrong
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		s, err := solver.NewSolver(puzzleInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if flags.All {
			solutionCount := 0
			for s.Solve() {
				iterations += s.Iterations()
				solutionCount++
				if solutionCount > flags.Limit && flags.Limit != 0 {
					globalLimit = true
					break
				}
				if !flags.CountsOnly {
					fmt.Printf("%s\n\n", format.Format(s.Solution(), flags.OutputFormat))
				}
			}
			// after all solutions of the current puzzle found
			totalSolutions += solutionCount
			if flags.CountsOnly {
				var count string
				if solutionCount > flags.Limit && flags.Limit != 0 {
					// indicate that we hit the limit, and hence the acutal number is higher
					count = fmt.Sprintf("%d (limit)", flags.Limit)
				} else {
					count = fmt.Sprintf("%d", solutionCount)
				}
				if flags.OutputInputPuzzle {
					fmt.Printf("%s: %s\n", format.Format(puzzleInput, "inline"), count)
				} else {
					fmt.Printf("%s\n", count)
				}
			}
		} else {
			if s.Solve() {
				iterations += s.Iterations()
				totalSolutions++
				fmt.Printf("%s\n", format.Format(s.Solution(), flags.OutputFormat))
			} else {
				fmt.Printf("No solution\n")
			}
		}
	}
	if flags.ShowStats {
		limit := ""
		if globalLimit {
			limit = " (limit)"
		}
		fmt.Printf("Total puzzles: %d\n", puzzleCount)
		fmt.Printf("Total solutions: %d%s\n", totalSolutions, limit)
		fmt.Printf("Total iterations: %d\n", iterations)
		fmt.Printf("Time taken: %s", time.Since(start))
	}
}
