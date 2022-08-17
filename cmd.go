package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AndrewSav/sudocoo/pkg/format"
)

type Flags struct {
	InputFile         string    // input can come from a file
	Input             string    // or form a string
	All               bool      // we want all solutions, not just the first one
	Limit             int       // we want that many first solutions of each puzzle
	CountsOnly        bool      // we want only solution counts, not soluctions themselves
	OutputInputPuzzle bool      // display puzzle along with its solution count
	OutputFormat      string    // how to print out a solution
	InputReader       io.Reader // we convert InputFile or Input to a uniform io.Reader
	ShowStats         bool      // display stats at the end of the program run
}

// this is so we could pring available output formats in usage help
func getAvailableFormats() string {
	const separator = ", "
	var sb strings.Builder
	s := []string{}
	for f := range format.GetKnownFormats() {
		s = append(s, f)
	}
	// since formats come from a map we have to sort it
	sort.Strings(s)
	for _, f := range s {
		fmt.Fprintf(&sb, "%s%s", f, separator)
	}
	result := sb.String()
	if len(result) > len(separator) {
		result = result[:len(result)-len(separator)]
	}
	return result
}

// see if user specified format is one of known formats
func validateFormat(s string) bool {
	for f := range format.GetKnownFormats() {
		if s == f {
			return true
		}
	}
	return false
}

func ParseArgs() Flags {
	var flags = Flags{}

	fs := flag.NewFlagSet("sudocoo", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Println("This is brute force solver for sudoku puzzles")
		fmt.Println("Based on code by Glenn Fowler of ATT http://gsf.cococlyde.org/")
		fmt.Println("Code archive: https://github.com/1to9only/ast-sudoku.2012-08-01")
		fmt.Printf("Usage: %s [FLAGS...]\n", filepath.Base(os.Args[0]))
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	fs.StringVar(&flags.InputFile, "f", "", "path to input file with puzzle(s). Only one of '-f' and '-i' can be specified")
	fs.StringVar(&flags.Input, "i", "", "puzzle input in inline format. You can specify a single asterisk '*' as the input to represent an empty puzzle. Only one of '-f' and '-i' can be specified")

	fs.BoolVar(&flags.All, "a", false, "find all solution, for each puzzle but no more than specified in the -l flag")
	fs.IntVar(&flags.Limit, "l", 1000, "the maximum number of solutions to find for each puzzle. 0 is no limit. Default: 1000. Only considered when '-a' is specified")

	fs.BoolVar(&flags.CountsOnly, "c", false, "do not print out the solutions, only solutions counts. Only considered when '-a' is specified")
	fs.BoolVar(&flags.OutputInputPuzzle, "p", false, "print puzzle intput in inline format along with each count. Only considered when '-c' is specified")

	fs.StringVar(&flags.OutputFormat, "v", "visual", fmt.Sprintf("output format for solutions: %s. Default: visual", getAvailableFormats()))
	fs.BoolVar(&flags.ShowStats, "s", false, "display total number of puzzles and solutions encountered and iterations taken at the end")

	fs.Parse(os.Args[1:])

	if fs.NArg() != 0 {
		fmt.Printf("want 0 arguments, have %d\n", fs.NArg())
		fs.Usage()
		os.Exit(2)
	}

	if flags.InputFile == "" && flags.Input == "" {
		fmt.Printf("you have to specify input with either -f or -i\n")
		fs.Usage()
		os.Exit(2)
	}

	if flags.InputFile != "" && flags.Input != "" {
		fmt.Printf("you have to specify either -f or -i, not both\n")
		fs.Usage()
		os.Exit(2)
	}

	if flags.InputFile != "" {
		file, err := os.Open(flags.InputFile)
		if err != nil {
			fmt.Printf("Error opening input file: %v\n", err)
			os.Exit(2)
		}
		flags.InputReader = file
	}

	if flags.Input != "" {
		if flags.Input == "*" {
			flags.InputReader = strings.NewReader(".................................................................................")
		} else {
			flags.InputReader = strings.NewReader(flags.Input)
		}

	}

	if !validateFormat(flags.OutputFormat) {
		fmt.Printf("invalid output format %s\n", flags.OutputFormat)
		fs.Usage()
		os.Exit(2)
	}

	return flags
}
