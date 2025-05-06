/*
Manages the data storage of the server.
*/
package data

import (
	"sync"

	"github.com/matetirpak/chessbot-playground-server/internal/game_logic"
)

type Game struct {
	Name         string
	ID           int32
	Password     string
	Started      bool
	HasWPlayer   bool
	WPlayerToken string
	HasBPlayer   bool
	BPlayerToken string
	Winner       string
	BoardData    []game_logic.BoardState
	Mu           sync.RWMutex
}

var GamesMap = make(map[int32]*Game)
var GamesMapMu sync.RWMutex

var NextFreeBoardID int32 = 1
var NextFreeBoardIDMu sync.RWMutex
