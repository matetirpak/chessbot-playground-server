/*
Unittest for the game_logic package.
*/
package game_logic

import (
	"sort"
	"testing"
)

func TestGenerateMoves(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', 'x', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'r', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', 'R', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'K', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		},
	}
	expectedMovesWhite := []*Move{
		{From: [2]int{1, 1}, To: [2]int{0, 0}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{0, 1}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{0, 2}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{1, 2}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{2, 2}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{2, 1}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{2, 0}, Color: 'w', Capture: false},
		{From: [2]int{1, 1}, To: [2]int{1, 0}, Color: 'w', Capture: false},

		{From: [2]int{3, 2}, To: [2]int{2, 2}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{1, 2}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{0, 2}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{3, 1}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{3, 0}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{4, 2}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{5, 2}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{6, 2}, Color: 'w', Capture: true},
		{From: [2]int{3, 2}, To: [2]int{3, 3}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{3, 4}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{3, 5}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{3, 6}, Color: 'w', Capture: false},
		{From: [2]int{3, 2}, To: [2]int{3, 7}, Color: 'w', Capture: false},
	}
	expectedMovesBlack := []*Move{
		{From: [2]int{6, 2}, To: [2]int{7, 0}, Color: 'b', Capture: false},
		{From: [2]int{6, 2}, To: [2]int{5, 0}, Color: 'b', Capture: false},
		{From: [2]int{6, 2}, To: [2]int{4, 1}, Color: 'b', Capture: false},
		{From: [2]int{6, 2}, To: [2]int{4, 3}, Color: 'b', Capture: false},
		{From: [2]int{6, 2}, To: [2]int{5, 4}, Color: 'b', Capture: false},
		{From: [2]int{6, 2}, To: [2]int{7, 4}, Color: 'b', Capture: false},

		{From: [2]int{4, 5}, To: [2]int{4, 4}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{4, 3}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{4, 2}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{4, 1}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{4, 0}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{3, 5}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{2, 5}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{1, 5}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{0, 5}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{4, 6}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{4, 7}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{5, 5}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{6, 5}, Color: 'b', Capture: false},
		{From: [2]int{4, 5}, To: [2]int{7, 5}, Color: 'b', Capture: false},
	}
	whiteMovesList := AllPossibleMoves('w', boardState, []rune{})
	if whiteMovesList == nil {
		t.Errorf("failed to generate moves for white")
	}
	whiteMoves := movesListToSlice(t, whiteMovesList)
	CompareMoves(t, whiteMoves, expectedMovesWhite)

	blackMovesList := AllPossibleMoves('b', boardState, []rune{})
	if whiteMovesList == nil {
		t.Errorf("failed to generate moves for black")
	}
	blackMoves := movesListToSlice(t, blackMovesList)
	CompareMoves(t, blackMoves, expectedMovesBlack)

}

// TestGenerateMovesForPiece tests the generateMovesForPiece function
func TestGenerateMovesForPiece(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{'x', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		},
	}

	// Define the expected moves for the King at position (0, 0)
	expectedMoves := []*Move{
		{From: [2]int{0, 0}, To: [2]int{1, 0}, Color: 'w', Capture: false},
		{From: [2]int{0, 0}, To: [2]int{1, 1}, Color: 'w', Capture: false},
		{From: [2]int{0, 0}, To: [2]int{0, 1}, Color: 'w', Capture: false},
	}

	var moves *Move
	GenerateMovesForPiece(0, 0, boardState, &moves)
	if moves == nil {
		t.Errorf("failed to generate moves for piece")
	}
	actualMoves := movesListToSlice(t, moves)

	CompareMoves(t, actualMoves, expectedMoves)

	boardState = &BoardState{
		Board: [8][8]rune{
			{'R', 'K', 'B', 'Q', 'X', 'B', 'K', 'R'},
			{'P', 'P', 'P', 'P', ' ', 'P', 'P', 'P'},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', 'P', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', 'k', ' ', ' '},
			{'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p'},
			{'r', 'k', 'b', 'q', 'x', 'b', ' ', 'r'},
		},
	}
	expectedMoves = []*Move{
		{From: [2]int{5, 5}, To: [2]int{4, 3}, Color: 'w', Capture: false},
		{From: [2]int{5, 5}, To: [2]int{3, 4}, Color: 'w', Capture: true},
		{From: [2]int{5, 5}, To: [2]int{3, 6}, Color: 'w', Capture: false},
		{From: [2]int{5, 5}, To: [2]int{4, 7}, Color: 'w', Capture: false},
		{From: [2]int{5, 5}, To: [2]int{7, 6}, Color: 'w', Capture: false},
	}
	moves = nil
	GenerateMovesForPiece(5, 5, boardState, &moves)
	if moves == nil {
		t.Errorf("failed to generate moves for piece")
	}
	actualMoves = movesListToSlice(t, moves)

	CompareMoves(t, actualMoves, expectedMoves)
}

// TestGenerateMovesForEmptySquare tests generateMovesForPiece with an empty square
func TestGenerateMovesForEmptySquare(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		},
	}

	var head *Move
	// Generate moves for the empty square at position (0, 0)
	GenerateMovesForPiece(0, 0, boardState, &head)

	// Check that no moves are generated
	if head != nil {
		t.Errorf("expected no moves for an empty square, but got: %+v", head)
	}
}

func TestGetDirectionDeltas(t *testing.T) {
	row := 3
	col := 3
	king_row := 5
	king_col := 5

	drow, dcol, err := getDirectionDeltas(row, col, king_row, king_col)
	if err != nil {
		t.Errorf("fail in GetDirectionDeltas: %s", err)
	}
	if drow != -1 || dcol != -1 {
		t.Errorf("both %d and %d should be -1. ", drow, dcol)
	}

	row = 2
	col = 3
	king_row = 5
	king_col = 0

	drow, dcol, err = getDirectionDeltas(row, col, king_row, king_col)
	if err != nil {
		t.Errorf("fail in GetDirectionDeltas: %s", err)
	}
	if drow != -1 || dcol != 1 {
		t.Errorf("%d should be -1 and %d should be 1. ", drow, dcol)
	}
}

func TestIsPinned(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{'b', Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, 'R', Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, 'X', Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, 'x', Empty, Empty, Empty, Empty, Empty},
		},
		WhiteKingPos: [2]int{7, 2},
		BlackKingPos: [2]int{3, 3},
	}
	drow, dcol, err := getDirectionDeltas(2, 2, 3, 3)
	if err != nil {
		t.Errorf("fail in GetDirectionDeltas: %s", err)
	}
	if drow != -1 || dcol != -1 {
		t.Errorf("both %d and %d should be -1. ", drow, dcol)
	}

	sol, err := isPinned(2, 2, boardState)
	if err != nil {
		t.Errorf("fail in IsPinned: %s", err)
	}
	if !sol {
		t.Errorf("solution should be true.")
	}

	boardState = &BoardState{
		Board: [8][8]rune{
			{'b', Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, 'R', Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, 'X', Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, 'x', Empty, Empty, Empty, Empty, Empty},
		},
		WhiteKingPos: [2]int{7, 2},
		BlackKingPos: [2]int{3, 2},
	}
	sol, err = isPinned(2, 2, boardState)
	if err != nil {
		t.Errorf("fail in IsPinned: %s", err)
	}
	if sol {
		t.Errorf("solution should be false.")
	}
}

func TestFieldAttacked(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{'b', Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
			{Empty, Empty, 'X', Empty, 'x', Empty, Empty, Empty},
		},
		WhiteKingPos: [2]int{7, 4},
		BlackKingPos: [2]int{7, 2},
	}
	sol, err := fieldAttacked(1, 1, 'w', boardState)
	if err != nil {
		t.Errorf("fail in FieldAttack: %s", err)
	}
	if !sol {
		t.Errorf("solution should be true.")
	}

	boardState = &BoardState{
		Board: [8][8]rune{
			{' ', 'r', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', 'B', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', 'X', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'x', ' ', ' ', ' ', ' ', ' '},
		},
		WhiteKingPos: [2]int{7, 2},
		BlackKingPos: [2]int{4, 1},
	}
	sol, err = isPinned(2, 1, boardState)
	if err != nil {
		t.Errorf("fail in IsPinned: %s", err)
	}
	if !sol {
		t.Errorf("solution should be true.")
	}

	sol, err = fieldAttacked(4, 3, 'b', boardState)
	if err != nil {
		t.Errorf("fail in FieldAttack: %s", err)
	}
	if !sol {
		t.Errorf("solution should be true.")
	}
}

func TestIsCheckmate(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{' ', 'X', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{'q', ' ', ' ', 'b', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', 'x', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		},
		WhiteKingPos: [2]int{4, 1},
		BlackKingPos: [2]int{0, 1},
	}
	checkmate, err := isCheckmatePlayer('b', *boardState)
	if err != nil {
		t.Errorf("fail in IsCheckmate: %s", err)
	}
	if !checkmate {
		t.Errorf("checkmate should be true")
	}
	winner, err := isCheckmate(*boardState)
	if err != nil {
		t.Errorf("fail in IsCheckmate: %s", err)
	}
	if winner != 'w' {
		t.Errorf("winner should be 'w'")
	}
}

func TestRemisPlayer(t *testing.T) {
	boardState := &BoardState{
		Board: [8][8]rune{
			{' ', 'X', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{'q', ' ', ' ', ' ', 'x', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'r', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		},
		WhiteKingPos: [2]int{2, 4},
		BlackKingPos: [2]int{0, 1},
	}
	blackRemis, err := isRemisPlayer('b', *boardState)
	if err != nil {
		t.Errorf("fail in IsRemisPlayer: %s", err)
	}
	if !blackRemis {
		t.Errorf("remis should be true")
	}
	remis, err := isRemis(*boardState)
	if err != nil {
		t.Errorf("fail in IsRemis: %s", err)
	}
	if !remis {
		t.Errorf("remis should be true")
	}
	boardState = &BoardState{
		Board: [8][8]rune{
			{' ', 'X', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{'q', ' ', ' ', ' ', 'x', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'r', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', 'R', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		},
		WhiteKingPos: [2]int{2, 4},
		BlackKingPos: [2]int{0, 1},
	}
	blackRemis, err = isRemisPlayer('b', *boardState)
	if err != nil {
		t.Errorf("fail in IsRemisPlayer: %s", err)
	}
	if blackRemis {
		t.Errorf("remis should be false")
	}
	remis, err = isRemis(*boardState)
	if err != nil {
		t.Errorf("fail in IsRemis: %s", err)
	}
	if remis {
		t.Errorf("remis should be false")
	}
}

func TestRemis(t *testing.T) {

}

func movesListToSlice(t *testing.T, moves *Move) []*Move {
	var ret []*Move
	ptmp := moves
	i := 0
	for ptmp != nil {
		i++
		if i == 10000 {
			t.Errorf("loop in movesListToSlice didn't terminate")
			return nil
		}
		ret = append(ret, ptmp)
		ptmp = ptmp.Next
	}
	return ret
}

func SortMoves(moves []*Move) {
	sort.SliceStable(moves, func(i, j int) bool {
		// Sorting first by from, then by to
		if moves[i].From != moves[j].From {
			return moves[i].From[0] < moves[j].From[0] || (moves[i].From[0] == moves[j].From[0] && moves[i].From[1] < moves[j].From[1])
		}
		return moves[i].To[0] < moves[j].To[0] || (moves[i].To[0] == moves[j].To[0] && moves[i].To[1] < moves[j].To[1])
	})
}

func CompareMoves(t *testing.T, actualMoves []*Move, expectedMoves []*Move) {
	// Sort the moves for comparison (if order doesn't matter)
	SortMoves(expectedMoves)
	SortMoves(actualMoves)

	// Manually compare the expected and actual moves
	if len(expectedMoves) != len(actualMoves) {
		t.Errorf("expected %d moves, but got %d moves", len(expectedMoves), len(actualMoves))
		return
	}

	// Compare each move
	for i := range expectedMoves {
		expectedMove := expectedMoves[i]
		actualMove := actualMoves[i]

		if expectedMove.From != actualMove.From ||
			expectedMove.To != actualMove.To ||
			expectedMove.Color != actualMove.Color ||
			expectedMove.Capture != actualMove.Capture {
			t.Errorf("move %d mismatch:\nExpected: %+v\nGot: %+v", i+1, expectedMove, actualMove)
		}
	}
}
