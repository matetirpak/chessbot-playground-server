/*
This module implements king functionalities to ensure
correct move validations and allow the detection
of terminal board states.
*/

package game_logic

import (
	"errors"
	"math"
)

func isPinned(row, col int, boardState *BoardState) (bool, error) {
	color, piece := getColorAndPiece(row, col, boardState.Board)
	if piece == 'x' {
		return false, nil
	}
	var king_row, king_col int
	if color == 'w' {
		king_row, king_col = boardState.WhiteKingPos[0], boardState.WhiteKingPos[1]
	} else {
		king_row, king_col = boardState.BlackKingPos[0], boardState.BlackKingPos[1]
	}
	if !isLinearCorrelated(row, col, king_row, king_col) {
		return false, nil
	}
	drow, dcol, err := getDirectionDeltas(row, col, king_row, king_col)
	if err != nil {
		return false, err
	}
	var diagonal bool = drow != 0 && dcol != 0
	var straight bool = drow == 0 || dcol == 0

	row_i, col_i := row, col
	i := 0
	for {
		i++
		if i == 10 {
			return false, errors.New("loop didn't terminate")
		}
		row_i += drow
		col_i += dcol
		if !isInBounds([2]int{row_i, col_i}) {
			break
		}
		target_color, target_piece := getColorAndPiece(row_i, col_i, boardState.Board)
		if target_piece == Empty {
			continue
		}
		if target_color == color {
			return false, nil
		}
		if diagonal && (target_piece == 'b' || target_piece == 'q') {
			return true, nil
		}
		if straight && (target_piece == 'r' || target_piece == 'q') {
			return true, nil
		}
	}
	return false, nil
}

func isLinearCorrelated(row1, col1 int, row2, col2 int) bool {
	return row1 == row2 || col1 == col2 || col1-col2 != 0 && math.Abs(float64((row1-row2)/(col1-col2))) == 1
}

func getDirectionDeltas(row, col int, king_row, king_col int) (int, int, error) {
	if !isLinearCorrelated(row, col, king_row, king_col) {
		return 0, 0, errors.New("can't generate direction deltas for uncorrelated fields")
	}

	// E or W
	if row == king_row {
		if col < king_col {
			// W
			return 0, -1, nil
		} else {
			// E
			return 0, 1, nil
		}
	}

	// N or S
	if col == king_col {
		if row < king_row {
			// N
			return -1, 0, nil
		} else {
			// S
			return 1, 0, nil
		}
	}

	// NE
	if row < king_row && col > king_col {
		return -1, 1, nil
	}

	// SE
	if row > king_row && col > king_col {
		return 1, 1, nil
	}

	// SW
	if row > king_row && col < king_col {
		return 1, -1, nil
	}

	// NW
	if row < king_row && col < king_col {
		return -1, -1, nil
	}
	return 0, 0, errors.New("delta direction couldn't be determined")
}

func fieldAttacked(row, col int, attackerColor rune, boardState *BoardState) (bool, error) {
	tmp := boardState.Board[row][col]
	if attackerColor == 'w' {
		boardState.Board[row][col] = 'X'
	} else {
		boardState.Board[row][col] = 'x'
	}
	var allMoves *Move = AllPossibleMoves(attackerColor, boardState, []rune{})
	boardState.Board[row][col] = tmp
	i := 0
	for allMoves != nil {
		if i == 1000 {
			return false, errors.New("error when traversing moves in FieldAttacked")
		}
		i++
		to_row, to_col := allMoves.To[0], allMoves.To[1]
		capture := allMoves.Capture
		allMoves = allMoves.Next

		if to_row != row || to_col != col {
			continue
		}
		if !capture {
			continue
		}
		return true, nil
	}
	return false, nil
}

func kingAttacked(color rune, boardState *BoardState) (bool, error) {
	king_row, king_col, enemy_color := getKingdataFromColor(color, boardState)
	attacked, err := fieldAttacked(king_row, king_col, enemy_color, boardState)
	return attacked, err
}

func kingMoveable(color rune, boardState *BoardState) (bool, error) {
	row, col, enemy_color := getKingdataFromColor(color, boardState)

	var moves *Move
	GenerateMovesForPiece(row, col, boardState, &moves)

	for moves != nil {
		to_row, to_col := moves.To[0], moves.To[1]
		moves = moves.Next
		attacked, err := fieldAttacked(to_row, to_col, enemy_color, boardState)
		if err != nil {
			return false, err
		}
		if !attacked {
			return true, nil
		}
	}

	return false, nil
}

// Checks if the given player was checkmated
func isCheckmatePlayer(color rune, boardState BoardState) (bool, error) {
	row, col, enemy_color := getKingdataFromColor(color, &boardState)

	// Check if king is indeed attacked
	attacked, err := fieldAttacked(row, col, enemy_color, &boardState)
	if err != nil {
		return false, err
	}
	if !attacked {
		return false, nil
	}

	// Check if king can move away
	dodge, err := kingMoveable(color, &boardState)
	if err != nil {
		return false, err
	}
	if dodge {
		return false, nil
	}

	// Check if another piece can block
	i := 0
	moves := AllPossibleMoves(color, &boardState, []rune{'x'})
	for moves != nil {
		i++
		if i >= 1000 {
			return false, errors.New("move loop didn't terminate")
		}
		from_row, from_col := moves.From[0], moves.From[1]
		to_row, to_col := moves.To[0], moves.To[1]
		moves = moves.Next
		tmp_from := boardState.Board[from_row][from_col]
		tmp_to := boardState.Board[to_row][to_col]
		boardState.Board[from_row][from_col] = Empty
		boardState.Board[to_row][to_col] = tmp_from

		attacked, err = fieldAttacked(row, col, enemy_color, &boardState)

		// Restore board
		boardState.Board[from_row][from_col] = tmp_from
		boardState.Board[to_row][to_col] = tmp_to

		if err != nil {
			return false, err
		}
		if !attacked {
			return false, nil
		}
	}
	return true, nil
}

// Returns the winning player if checkmate
func isCheckmate(boardState BoardState) (rune, error) {
	whiteCheckmated, err := isCheckmatePlayer('w', boardState)
	if err != nil {
		return 'n', err
	}
	if whiteCheckmated {
		return 'b', nil
	}
	blackCheckmated, err := isCheckmatePlayer('b', boardState)
	if err != nil {
		return 'n', err
	}
	if blackCheckmated {
		return 'w', nil
	}
	return 'n', nil
}

// Checks if the player can't move
func isRemisPlayer(color rune, boardState BoardState) (bool, error) {
	row, col, enemy_color := getKingdataFromColor(color, &boardState)
	attacked, err := fieldAttacked(row, col, enemy_color, &boardState)
	if err != nil {
		return false, err
	}
	if attacked {
		return false, nil
	}

	moveable, err := kingMoveable(color, &boardState)
	if err != nil {
		return false, err
	}
	if moveable {
		return false, nil
	}

	moves := AllPossibleMoves(color, &boardState, []rune{'x'})
	if moves != nil {
		return false, nil
	}
	return true, nil
}

func isRemis(boardState BoardState) (bool, error) {
	whiteRemis, err := isRemisPlayer('w', boardState)
	if err != nil {
		return false, err
	}
	if whiteRemis {
		return true, nil
	}
	blackRemis, err := isRemisPlayer('b', boardState)
	if err != nil {
		return false, err
	}
	if blackRemis {
		return true, nil
	}
	return false, nil
}

func getKingdataFromColor(color rune, boardState *BoardState) (int, int, rune) {
	var row, col int
	var enemy_color rune
	if color == 'w' {
		row, col = boardState.WhiteKingPos[0], boardState.WhiteKingPos[1]
		enemy_color = 'b'
	} else {
		row, col = boardState.BlackKingPos[0], boardState.BlackKingPos[1]
		enemy_color = 'w'
	}
	return row, col, enemy_color
}

func FilterInvalidMoves(move *Move, bstate *BoardState) (*Move, error) {
	var retHead *Move
	movep := &retHead

	for move != nil {
		nextBState := MakeMove(move, *bstate, false)
		check, err := kingAttacked(move.Color, &nextBState)
		if err != nil {
			return retHead, err
		}
		if check {
			move = move.Next
			continue
		}
		tmpMove := *move
		tmpMove.Next = nil
		*movep = &tmpMove
		movep = &(tmpMove.Next)
		move = move.Next
	}

	return retHead, nil
}
