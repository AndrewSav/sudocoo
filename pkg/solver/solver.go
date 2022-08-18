package solver

import (
	"fmt"
)

// Algorithm outline: find the cell with fewest candidates. Put one of the candidates in the cell.
// Repeat until either all cells are filled (solution found) or there is a cell with no candidates
// (solution no longer possible). In the latter case try another candidate for the current cell,
// if all candidates already tried, empty the cell in try another candidate for the previous cell, etc.

// Details: Keep a list of cells we filled with the candidates starting with an empty list.
// Attached to each cell keep a list of candidates to try for the cells in the list.
// We will call the last cell in the list of cells "the current (search) cell".
// Step 1: Find next cell to fill. If there are no empty cells left we've found a solution.
// Record the solution, and to continue search backtrack: go to Step 2.
// Otherwise, if there is a cell with no possible candidates, backtrack: go to Step 2.
// Otherwise find a cell that has fewest possible candidates, in case of a tie choose any.
// Put this cell at the end of the list of cells mentioned above. Initialize the attached
// candidates list with all possible candidates for this cell, fill the cell: go to Step 3.
// Step 2: Backtrack. If the current cell does not have any candidates in it's attached list,
// empty the cell and remove it from the list. Repeat until the current cell has candidate(s).
// If no cells left in the list, there no (more) solutions. Otherwise, fill the cell: go to Step 3.
// Step 3: Fill the cell. Remove a candidate from the candidates list attached to the current cell
// and put the candidate into the cell. Find next cell to fill, go to Step 1.

// Implementation note: during step 2 when we remove a cell from the grid and during step 3
// when we add cell to the grid the possible candidates for other cells change. In order to speed up
// calculating those on each iteration, we keep a data structure ('candidates' below) that we
// update each time we add or remove a cell. It allows us to calculate all possible candidates
// for a given cell in 3 lookups and 2 and operations ('getCellCandidates'), which is relatively fast.

const sudokuSize = 9

// Used to determine which box a particular cell identified by its row and column belongs to
var boxLookup = [sudokuSize][sudokuSize]int{
	{0, 0, 0, 1, 1, 1, 2, 2, 2},
	{0, 0, 0, 1, 1, 1, 2, 2, 2},
	{0, 0, 0, 1, 1, 1, 2, 2, 2},
	{3, 3, 3, 4, 4, 4, 5, 5, 5},
	{3, 3, 3, 4, 4, 4, 5, 5, 5},
	{3, 3, 3, 4, 4, 4, 5, 5, 5},
	{6, 6, 6, 7, 7, 7, 8, 8, 8},
	{6, 6, 6, 7, 7, 7, 8, 8, 8},
	{6, 6, 6, 7, 7, 7, 8, 8, 8},
}

// A candidate for a cell, row, column or box is any number
// that can go into that cell, row, column or box without conflicting with
// other numbers that are already in the grid according to sudoku rules.
// A cell candidate is an intersection of candidates of cell's row, column and box

// This struct keeps bit mask of possible candidates for each row, column and box
// only the right 9 bits are used of each int
// 0b100000100 means that 9 and 3 are possible candidates
// and the remaining numbers are eliminated because
// they are already present in this row, column or box
type candidates struct {
	row    [sudokuSize]int
	column [sudokuSize]int
	box    [sudokuSize]int
}

// This is how initial candidates start - all nine are possible
const initialCandidatesMask = 0b111111111

// This method is called when a number is added or remove to/from a cell in the solution
// x and y represent column and row respectively
// bit represents the number, e.g for 9 it will be 0b100000000 = 256,
// for 3 it will be 0b000000100 = 4, etc
func (c *candidates) flipBit(x, y, bit int) {
	c.row[y] ^= bit
	c.column[x] ^= bit
	c.box[boxLookup[y][x]] ^= bit
}

// This is a version of flipBit which is called during the puzzle initialization.
// Since all candidate bits start as '1' if after the flip the candidate is not
// zero, it means that the passed number already appeared in the row, column or box
// and hence the input is invalid
func (c *candidates) flipBitWithCheck(x, y, bit int) bool {
	c.flipBit(x, y, bit)
	return (c.row[y]&bit == 0) && (c.column[x]&bit == 0) && (c.box[boxLookup[y][x]]&bit == 0)
}

// For a given cell return all possible candidates, intersecting
// row, column and box candidates
func (c *candidates) getCellCandidates(x, y int) int {
	return c.row[y] & c.column[x] & c.box[boxLookup[y][x]]
}

// This represents a cell position in the sudoku grid
type coordinates struct {
	row    int
	column int
}

// Integers in 'cells' and 'cellCandidates' are bit fields
// 'cells' elements always have a single bit set - corresponding to the number in the cell
// or none if the cell is empty
// 'cellCandidates' will have as many bits set as there are candidates remaining to try
// In both case only nine right bits are used
type Solver struct {
	globalCandidates  candidates                  // candidates for each row, column and box
	cells             [sudokuSize][sudokuSize]int // sudoku cells, empty cells are zeroes
	cellSearchSpace   []coordinates               // list of empty cells that we are trying to fill to find solutions
	currentSearchCell int                         // the index of the current cell in the cellSearchSpace
	cellCandidates    [sudokuSize][sudokuSize]int // candidates for each cell to still try
	lastSolution      [sudokuSize][sudokuSize]int // copy of .cells as of last found solution
	done              bool                        // indicator that the solver has finished
	haveSolution      bool                        // indicator the .lastSolution contains a solution
	iterations        int                         // current iteration number for statistics purposes
}

// Flips the candidate bits for the current search cell, adding or removing the number in the current search cell to/from
// the candidate lists
func (s *Solver) flip() {
	s.globalCandidates.flipBit(s.cellSearchSpace[s.currentSearchCell].column, s.cellSearchSpace[s.currentSearchCell].row, s.getCurrentCell())
}

// Returns the remaining to try candidates for the current cell
func (s *Solver) getCurrentCellCandidates() int {
	return s.cellCandidates[s.cellSearchSpace[s.currentSearchCell].row][s.cellSearchSpace[s.currentSearchCell].column]
}

// Sets the remaining to try candidates for the current cell
func (s *Solver) setCurrentCellCandidates(c int) {
	s.cellCandidates[s.cellSearchSpace[s.currentSearchCell].row][s.cellSearchSpace[s.currentSearchCell].column] = c
}

// Returns the current cell value
func (s *Solver) getCurrentCell() int {
	return s.cells[s.cellSearchSpace[s.currentSearchCell].row][s.cellSearchSpace[s.currentSearchCell].column]
}

// Sets the current cell value
func (s *Solver) setCurrentCell(v int) {
	s.cells[s.cellSearchSpace[s.currentSearchCell].row][s.cellSearchSpace[s.currentSearchCell].column] = v
}

// This allows us converting numbers from Solver.cells to actual numbers, e.g. for solution results
var bitToNumber = [1 << sudokuSize]int{
	0:      0,
	1 << 0: 1,
	1 << 1: 2,
	1 << 2: 3,
	1 << 3: 4,
	1 << 4: 5,
	1 << 5: 6,
	1 << 6: 7,
	1 << 7: 8,
	1 << 8: 9,
}

// This lookup is used for getting the next candidate to try from a candidates bit mask
var leftmostBitLookup [1 << sudokuSize]int

// This lookup is used to get the number of remaining candidates from a candidates bit mask
var bitCount [1 << sudokuSize]int

// This is the initial .globalCandidates value (all rows, cells and boxes have all possible candidates) used in initialization of the Solver object
var initialCandidates candidates

func init() {
	// Populate leftmostBitLookup and bitCount above
	for i := 0; i < 1<<sudokuSize; i++ {
		count := 0
		for j := 1; j < 1<<sudokuSize; j <<= 1 { // for each bit in i
			if i&j != 0 { // if this bit is set in i
				count++
				leftmostBitLookup[i] = j // this will eventually gets overwritten in the inner loop by the leftmost value
			}
		}
		bitCount[i] = count
	}
	// Populate initialCandidates above
	for i := 0; i < sudokuSize; i++ {
		initialCandidates.row[i] = initialCandidatesMask
		initialCandidates.column[i] = initialCandidatesMask
		initialCandidates.box[i] = initialCandidatesMask
	}
}

// Create a new solver from 9x9 integer array of sudoku input
// Returns error when the input array is inconsistent (same number in a row, column or box)
func NewSolver(s [sudokuSize][sudokuSize]int) (*Solver, error) {
	sudoku := Solver{globalCandidates: initialCandidates, currentSearchCell: -1}
	for y, row := range sudoku.cells {
		for x := range row {
			digit := s[y][x]
			if digit > 0 {
				digit = 1 << (digit - 1)
				// Adjust candidates table to account for this non-empty cell
				if !sudoku.globalCandidates.flipBitWithCheck(x, y, digit) {
					return nil, fmt.Errorf("invalid (inconsistent) puzzle input")
				}
			} else {
				// Add this empty cell into the search space
				sudoku.cellSearchSpace = append(sudoku.cellSearchSpace, coordinates{y, x})
			}
			// put the cell in the grid
			sudoku.cells[y][x] = digit
		}
	}
	return &sudoku, nil
}

// Call this after a prior call to .Solve() returned true
func (s *Solver) Solution() (result [sudokuSize][sudokuSize]int) {
	if !s.haveSolution {
		panic("Solution is called before Solve returned true")
	}
	for y, row := range s.lastSolution {
		for x := range row {
			result[y][x] = bitToNumber[s.lastSolution[y][x]]
		}
	}
	return
}

// Returns the number of iterations performed for statistical purposes
func (s *Solver) Iterations() int {
	return s.iterations
}

// Find next empty cell to try. Returns true if no more cells to try, and thus
// we found a solution. There are three possible outcomes:
//  1. As above, no more cells to try, all are filled, we always fill according to the rules
//     so this is a solution
//  2. There is a cell that does not allow any candidates. In this case we are staying on the current
//     cell so next candidate of this cell (if any) could be tried
//  3. Otherwise find an empty cell with the least possible number of candidates, and set it current
func searchNextCellToTry(s *Solver) bool {
	var (
		// The maximum number of candidates for any cell is 9, so we start search with 10,
		// indicating that any possible number is an improvement
		fewestCandidatesCount = 10
		// Index of the "best" found so far cell in cellSearchSpace
		indexFound = -1
		// Cell candidates of the "best" found so far cell in cellSearchSpace
		cellCandidates = 0
	)
	// If we ran out of empty cells we have a solution
	if s.currentSearchCell == len(s.cellSearchSpace)-1 {
		return true
	}
	// All the empty cells has higher index than the current cell in cellSearchSpace
	for i := s.currentSearchCell + 1; i < len(s.cellSearchSpace); i++ {
		// Get cell candidates for the cell
		cc := s.globalCandidates.getCellCandidates(s.cellSearchSpace[i].column, s.cellSearchSpace[i].row)
		// Get the number of candidates
		bc := bitCount[cc]
		// If no candidates, no point searching further,
		// let's return so we can try other remaining candidates
		// for the current search sell if any, or backtrack if none
		if bc == 0 {
			return false
		}
		if fewestCandidatesCount > bc {
			// We found a cell with less candidates than previous best
			// lets store its cell candidates and index, and update current best candidates count
			cellCandidates = cc
			indexFound = i
			fewestCandidatesCount = bc
			// if we have only 1 candidate, we cannot improve that any further
			// so we might as well stop searching
			if fewestCandidatesCount == 1 {
				break
			}
		}
	}
	// Swap the found cell with the current one, in the cellSearchSpace since we going to occupy it shortly with a candidate
	// and we want empty cells to remain at the end as we rely on this in the code above
	s.currentSearchCell++
	if indexFound != s.currentSearchCell {
		s.cellSearchSpace[indexFound], s.cellSearchSpace[s.currentSearchCell] = s.cellSearchSpace[s.currentSearchCell], s.cellSearchSpace[indexFound]
	}
	s.setCurrentCellCandidates(cellCandidates)
	return false
}

// This backtracks to the previous cell that still
// has remaining candidates. Return 0 if we backtracked
// to the start, otherwise returns new current cell candidates
func (s *Solver) backtrack() int {
	// If we are here the current cell has no candidates left
	for {
		// Restore global candidates table by removing
		// the number in the current cell
		s.flip()
		// Make previous cell current
		s.currentSearchCell--
		// If we are back to start we finished the search
		if s.currentSearchCell == -1 {
			return 0
		}
		// Get candidates for now current cell
		lcc := s.getCurrentCellCandidates()
		// If we have candidates return them,
		// but restore global candidates table once again
		// since we about to replace that cell current number
		if lcc != 0 {
			s.flip()
			return lcc
		}
	}
}

// Call this to find next solution. Returns false when no more solutions
// and true, when a solution is found. After true is returned call
// .Solution() to get last solution
func (s *Solver) Solve() bool {
	// Sometimes we discover that we completed the full search
	// and cannot backtrack any further on the same iteration
	// when we find the last solution, but since Solve() returns
	// true, the client is likely to call it again, we need
	// to return false in this case
	if s.done {
		return false
	}
	for {
		s.iterations++ // in theory this can overflow, in practice it would take too long
		// Find next cell to try
		haveSolution := searchNextCellToTry(s)
		// If all cells are filled it's a solution
		if haveSolution {
			s.haveSolution = true    // so .Solution() could panic if there is no solution yey
			s.lastSolution = s.cells // we'll move on soon, so store it for .Solution() to return
		}
		// Get candidates for the selected cell
		lcc := s.getCurrentCellCandidates()
		// If no candidates, we need to backtrack
		if lcc == 0 {
			lcc = s.backtrack()
			// If we cannot backtrack any further the search if finished
			if lcc == 0 {
				s.done = true
				// We might also have found a solution on the same iteration
				// if so, indicate it to the caller
				return haveSolution
			}
		}
		// Get next candidate
		candidate := leftmostBitLookup[lcc]
		// Remove the candidate from the cell's candidates list
		s.setCurrentCellCandidates(lcc ^ candidate)
		// Write the candidate to the cell
		s.setCurrentCell(candidate)
		// Update global candidates table, to indicate that this number is no longer candidate
		// for the respective row, column and box
		s.flip()
		// if we found a solution earlier, indicate it to the caller
		if haveSolution {
			return true
		}
	}
}
