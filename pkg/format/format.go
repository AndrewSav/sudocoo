package format

import (
	"fmt"
	"strings"
)

const sudokuSize = 9

type FormatTemplate struct {
	Name                   string
	Description            string
	Header                 string
	ColumnSeparator        string
	RowSeparator           string
	VerticalBoxSeparator   string
	HorizontalBoxSeparator string
	Footer                 string
	Empty                  string // Empty cell character
	ColumnPrefix           string
	ColumnSuffix           string
}

// These formats come from here: https://github.com/1to9only/ast-sudoku.2012-08-01/blob/master/src/cmd/sudoku/sudocoo.rt
var formats = map[string]FormatTemplate{
	"inline": {
		Name:                   "inline",
		Description:            "Single line, Empty cells are dots",
		Header:                 "",
		ColumnSeparator:        "",
		RowSeparator:           "",
		VerticalBoxSeparator:   "",
		HorizontalBoxSeparator: "",
		Footer:                 "",
		Empty:                  ".",
		ColumnPrefix:           "",
		ColumnSuffix:           "",
	},
	"zeroes": {
		Name:                   "zeroes",
		Description:            "Single line, Empty cells are zeroes",
		Header:                 "",
		ColumnSeparator:        "",
		RowSeparator:           "",
		VerticalBoxSeparator:   "",
		HorizontalBoxSeparator: "",
		Footer:                 "",
		Empty:                  "0",
		ColumnPrefix:           "",
		ColumnSuffix:           "",
	},
	"visual": {
		Name:                   "visual",
		Description:            "Sudoku programmer's forum post format",
		Header:                 "",
		ColumnSeparator:        " ",
		RowSeparator:           "\n",
		VerticalBoxSeparator:   "|",
		HorizontalBoxSeparator: "---------------------\n",
		Footer:                 "",
		Empty:                  ".",
		ColumnPrefix:           "",
		ColumnSuffix:           "",
	},
	"vbforums": {
		Name:                   "vbforums",
		Description:            "VBForums Contest format (*.msk;*.sol)",
		Header:                 "",
		ColumnSeparator:        "",
		RowSeparator:           "\n",
		VerticalBoxSeparator:   "",
		HorizontalBoxSeparator: "",
		Footer:                 "",
		Empty:                  ".",
		ColumnPrefix:           "",
		ColumnSuffix:           "",
	},
	"sadman": {
		Name:                   "sadman",
		Description:            "SadMan Software Sudoku format (*.sdk)",
		Header:                 "[Puzzle]\n",
		ColumnSeparator:        "",
		RowSeparator:           "\n",
		VerticalBoxSeparator:   "",
		HorizontalBoxSeparator: "",
		Footer:                 "",
		Empty:                  ".",
		ColumnPrefix:           "",
		ColumnSuffix:           "",
	},
	"simple": {
		Name:                   "simple",
		Description:            "Simple Sudoku format (*.ss)",
		Header:                 "*-----------*\n",
		ColumnSeparator:        "",
		RowSeparator:           "\n",
		VerticalBoxSeparator:   "|",
		HorizontalBoxSeparator: "|---+---+---|\n",
		Footer:                 "\n*-----------*",
		Empty:                  ".",
		ColumnPrefix:           "|",
		ColumnSuffix:           "|",
	},
	"solver": {
		Name:                   "solver",
		Description:            "SuDoku Solver format (*.spf)",
		Header:                 "",
		ColumnSeparator:        "",
		RowSeparator:           "\n",
		VerticalBoxSeparator:   "|",
		HorizontalBoxSeparator: "---+---+---\n",
		Footer:                 "",
		Empty:                  ".",
		ColumnPrefix:           "",
		ColumnSuffix:           "",
	},
}

func FormatFromTemplate(puzzle [sudokuSize][sudokuSize]int, format FormatTemplate) string {
	var sb strings.Builder
	if format.Header != "" {
		fmt.Fprintf(&sb, "%s", format.Header)
	}
	for y := 0; y < sudokuSize; y++ {
		if format.ColumnPrefix != "" {
			fmt.Fprintf(&sb, "%s", format.ColumnPrefix)
		}
		for x := 0; x < sudokuSize; x++ {
			digit := fmt.Sprintf("%d", puzzle[y][x])
			if digit == "0" {
				fmt.Fprintf(&sb, format.Empty)
			} else {
				fmt.Fprintf(&sb, "%s", digit)
			}
			if format.ColumnSeparator != "" && x != sudokuSize-1 {
				fmt.Fprintf(&sb, "%s", format.ColumnSeparator)
			}
			if format.VerticalBoxSeparator != "" && x != sudokuSize-1 && x != 0 && x%3 == 2 {
				fmt.Fprintf(&sb, "%s", format.VerticalBoxSeparator)
				if format.ColumnSeparator != "" {
					fmt.Fprintf(&sb, "%s", format.ColumnSeparator)
				}
			}
		}
		if format.ColumnSuffix != "" {
			fmt.Fprintf(&sb, "%s", format.ColumnSuffix)
		}
		if format.RowSeparator != "" && y != sudokuSize-1 {
			fmt.Fprintf(&sb, "%s", format.RowSeparator)
		}
		if format.HorizontalBoxSeparator != "" && y != sudokuSize-1 && y != 0 && y%3 == 2 {
			fmt.Fprintf(&sb, "%s", format.HorizontalBoxSeparator)
		}
	}
	if format.Footer != "" {
		fmt.Fprintf(&sb, "%s", format.Footer)
	}
	return sb.String()
}

func Format(puzzle [sudokuSize][sudokuSize]int, formatName string) string {
	format, ok := formats[formatName]
	if !ok {
		panic(fmt.Sprintf("Unknown format '%s'", formatName))
	} else {
		return FormatFromTemplate(puzzle, format)
	}
}

func GetKnownFormats() map[string]FormatTemplate {
	result := make(map[string]FormatTemplate)
	for k, v := range formats {
		result[k] = v
	}
	return result
}
