package parser

import (
	"bufio"
	"fmt"
	"io"
)

const sudokuSize = 9

// Maps input characters to solver digits
var runeLookup = map[string]int{
	".": 0,
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
}

// Reads next puzzle input from bufio.Scanner,
// returns io.EOF when no more input,
// scanner needs to be created by parser.CreateInputScanner
func ReadNextPuzzleInput(s *bufio.Scanner) (result [sudokuSize][sudokuSize]int, err error) {
	for y := 0; y < sudokuSize; y++ {
		for x := 0; x < sudokuSize; x++ {
			hasDigit := false
			for s.Scan() {
				digit, ok := runeLookup[s.Text()]
				if !ok {
					continue
				}
				hasDigit = true
				result[y][x] = digit
				break
			}
			if !hasDigit {
				if x == 0 && y == 0 {
					err = io.EOF
					return
				} else {
					err = fmt.Errorf("Not enough valid sudoku characters ('.',0-9) in the input")
					return
				}
			}
		}
	}
	return
}

// Prepares runes scanner that parser.ReadNextPuzzleInput expects
func CreateInputScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	return scanner
}
