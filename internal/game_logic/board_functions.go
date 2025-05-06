/*
Board functionalities and helper functions used to extract
board information are implemented here.
*/

package game_logic

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// BoardState represents the chessboard and other game infos.
//
// Winner:
//   - 'n': none
//   - 'r': remis
//   - 'w': white
//   - 'b': black
//
// EnPassant:
//   - If a pawn double moves, contains coordinates of the target field.
//   - Otherwise: {-1, -1}
//
// TurnColor:
//   - 'w': white
//   - 'b': black
type BoardState struct {
	Board          [8][8]rune `json:"board"`
	LastMove       string     `json:"lastmove"`
	WhiteKingPos   [2]int     `json:"whitekingpos"`
	BlackKingPos   [2]int     `json:"blackkingpos"`
	WhiteKingMoved bool       `json:"whitekingmoved"`
	BlackKingMoved bool       `json:"blackkingmoved"`
	Winner         string     `json:"winner"`
	TurnColor      string     `json:"turncolor"`
	EnPassant      [2]int     `json:"enpassant"`
}

// Constructs the standard starting board.
func InitializeBoard(boardStates *[]BoardState) {
	var bstate BoardState
	bstate.LastMove = ""
	bstate.WhiteKingPos = [2]int{7, 4}
	bstate.BlackKingPos = [2]int{0, 4}
	bstate.WhiteKingMoved = false
	bstate.BlackKingMoved = false
	bstate.Winner = "n"
	bstate.EnPassant = [2]int{-1, -1}
	bstate.TurnColor = "n"

	board := &bstate.Board

	// Initialize black pieces
	board[0][0] = 'R'
	board[0][1] = 'K'
	board[0][2] = 'B'
	board[0][3] = 'Q'
	board[0][4] = 'X'
	board[0][5] = 'B'
	board[0][6] = 'K'
	board[0][7] = 'R'
	for i := 0; i < 8; i++ {
		board[1][i] = 'P'
	}

	// Initialize empty fields
	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			board[i][j] = Empty
		}
	}

	// Initialize white pieces
	for i := 0; i < 8; i++ {
		board[6][i] = 'p'
	}
	board[7][0] = 'r'
	board[7][1] = 'k'
	board[7][2] = 'b'
	board[7][3] = 'q'
	board[7][4] = 'x'
	board[7][5] = 'b'
	board[7][6] = 'k'
	board[7][7] = 'r'
	*boardStates = append(*boardStates, bstate)
}

// Returns the color and piece given a position on the board.
func getColorAndPiece(row int, col int, board [8][8]rune) (color rune, piece rune) {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		// Field out of bounds.
		// Error is not raised for simplicity.
		return 'n', Empty
	}

	// Extract piece
	p := board[row][col]
	if p == Empty {
		return 'n', Empty
	}

	// Return color and piece in lowercase
	if p >= 'a' && p <= 'z' {
		return 'w', rune(strings.ToLower(string(p))[0])
	}
	return 'b', rune(strings.ToLower(string(p))[0])
}

// Checks whether a position is within the board's bounds.
func isInBounds(pos [2]int) bool {
	return pos[0] >= 0 && pos[0] < 8 && pos[1] >= 0 && pos[1] < 8
}

// Applies a specified move to the board.
func MakeMove(move *Move, bstate BoardState, realMove bool) BoardState {
	fromRow, fromCol := move.From[0], move.From[1]
	toRow, toCol := move.To[0], move.To[1]
	fromColor, fromPiece := getColorAndPiece(fromRow, fromCol, bstate.Board)

	var newBstate BoardState
	data, _ := json.Marshal(bstate)
	json.Unmarshal(data, &newBstate)

	// Update king positions
	if fromPiece == 'x' {
		if fromColor == 'w' {
			newBstate.WhiteKingMoved = true
			newBstate.WhiteKingPos = [2]int{toRow, toCol}
		}
		if fromColor == 'b' {
			newBstate.BlackKingMoved = true
			newBstate.BlackKingPos = [2]int{toRow, toCol}
		}
	}

	// Update en passant
	newBstate.EnPassant = [2]int{-1, -1}
	if fromPiece == 'p' {
		if math.Abs(float64(fromRow-toRow)) == 2 {
			newBstate.EnPassant = [2]int{toRow, toCol}
		}
	}

	// Update last move
	newBstate.LastMove = MoveToString(move)

	// Update board
	newBstate.Board[toRow][toCol] = newBstate.Board[fromRow][fromCol]
	newBstate.Board[fromRow][fromCol] = Empty

	// Update player
	if fromColor == 'w' {
		newBstate.TurnColor = "b"
	} else {
		newBstate.TurnColor = "w"
	}

	if realMove {
		enemyColor := 'w'
		if fromColor == 'w' {
			enemyColor = 'b'
		}
		checkmate, err := isCheckmatePlayer(enemyColor, newBstate)
		if err != nil {
			fmt.Println("Error occured when checking for checkmate.")
		}
		if checkmate {
			newBstate.Winner = string(fromColor)
			newBstate.TurnColor = "n"
		}
		remis, err := isRemisPlayer(enemyColor, newBstate)
		if err != nil {
			fmt.Println("Error occured when checking for remis.")
		}
		if remis {
			newBstate.Winner = "r"
			newBstate.TurnColor = "n"
		}
	}
	return newBstate
}
