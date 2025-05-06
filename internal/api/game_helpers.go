/*
Helper functions for game data conversion.
*/

package api

import (
	gl "github.com/matetirpak/chessbot-playground-server/internal/game_logic"
)

func llToArray(move *gl.Move) []RespMove {
	var moves []RespMove
	for current := move; current != nil; current = current.Next {
		apiMove := RespMove{
			From:    current.From,
			To:      current.To,
			Capture: current.Capture,
		}
		moves = append(moves, apiMove)
	}
	return moves
}
