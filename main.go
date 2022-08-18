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
	var (
		totalSolutions = 0
		globalLimit    = false
		puzzleCount    = 0
		iterations     = 0
		start          = time.Now()
	)

	for ; ; puzzleCount++ {
		puzzleInput, err := parser.ReadNextPuzzleInput(scanner)
		// If this is the first puzzle and there is no puzzle,
		// then it's a error, otherwise we processed all puzzles
		if errors.Is(err, io.EOF) && puzzleCount != 0 {
			break
		}
		// We exit on these errors because the format is realy loose
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
		if flags.DontSolve {
			fmt.Printf("%s\n", format.Format(puzzleInput, flags.OutputFormat))
			if flags.NewLineAfterEachPuzzle {
				fmt.Println()
			}
		} else {
			if flags.All {
				solutionCount := 0
				for s.Solve() {
					iterations += s.Iterations()
					solutionCount++
					if solutionCount > flags.Limit && flags.Limit != 0 {
						solutionCount--
						globalLimit = true
						break
					}
					if !flags.CountsOnly && !(flags.ShowStats && flags.Quiet) {
						fmt.Printf("%s\n", format.Format(s.Solution(), flags.OutputFormat))
						if flags.NewLineAfterEachPuzzle {
							fmt.Println()
						}
					}
				}
				// After all solutions of the current puzzle found
				totalSolutions += solutionCount
				if flags.CountsOnly && !(flags.ShowStats && flags.Quiet) {
					var count string
					if solutionCount > flags.Limit && flags.Limit != 0 {
						// Indicate that we hit the limit, and hence the acutal number is higher
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
					if !(flags.ShowStats && flags.Quiet) {
						fmt.Printf("%s\n", format.Format(s.Solution(), flags.OutputFormat))
						if flags.NewLineAfterEachPuzzle {
							fmt.Println()
						}
					}
				} else {
					fmt.Printf("No solution\n")
				}
			}
		}
	}
	if flags.ShowStats {
		limit := ""
		if globalLimit {
			// Indicate that we hit the limit, and hence the acutal number is higher
			limit = " (limit)"
		}
		fmt.Printf("Total puzzles: %d\n", puzzleCount)
		fmt.Printf("Total solutions: %d%s\n", totalSolutions, limit)
		fmt.Printf("Total iterations: %d\n", iterations)
		fmt.Printf("Time taken: %s", time.Since(start))
	}
}
