/*
Helper functions for API communication and game setup.
*/
package api

import (
	"net/http"

	"github.com/google/uuid"

	data "github.com/matetirpak/chessbot-playground-server/internal/data"
	"github.com/matetirpak/chessbot-playground-server/internal/game_logic"
)

// Verifies whether a user has access to a session.
func verifyGameAccess(w http.ResponseWriter, id int32, password string) bool {
	data.GamesMapMu.RLock()
	game, exists := data.GamesMap[id]
	data.GamesMapMu.RUnlock()

	if !exists {
		http.Error(w, "Board not found", http.StatusNotFound)
		return false
	}

	game.Mu.RLock()
	gamePassword := game.Password
	game.Mu.RUnlock()

	if password != gamePassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return false
	}
	return true
}

// Verifies whether a user is registered as a player and has access to a session.
func verifyBoardAccess(w http.ResponseWriter, game *data.Game, color string, token string) bool {
	// Both ip and token have to match with the database, session password is not required.
	game.Mu.RLock()
	hasWPlayer := game.HasWPlayer
	wToken := game.WPlayerToken
	hasBPlayer := game.HasBPlayer
	bToken := game.BPlayerToken
	game.Mu.RUnlock()

	switch color {
	case "w":
		if !hasWPlayer {
			http.Error(w, "White player does not exist.", http.StatusNotFound)
			return false
		}
		if wToken != token {
			http.Error(w, "Token is invalid.", http.StatusUnauthorized)
			return false
		}
		return true

	case "b":
		if !hasBPlayer {
			http.Error(w, "Black player does not exist.", http.StatusNotFound)
			return false
		}
		if bToken != token {
			http.Error(w, "Token is invalid.", http.StatusUnauthorized)
			return false
		}
		return true

	default:
		http.Error(w, "Invalid color. Enter 'w' or 'b'", http.StatusBadRequest)
		return false
	}
}

func generateToken() string {
	return uuid.New().String()
}

func initializeNewGame(name string) *data.Game {
	var game data.Game

	data.NextFreeBoardIDMu.Lock()
	game.ID = data.NextFreeBoardID
	data.NextFreeBoardID++
	data.NextFreeBoardIDMu.Unlock()

	game.Name = name
	game.Password = generateToken()
	game.HasWPlayer = false
	game.HasBPlayer = false
	game.Winner = "n"
	game_logic.InitializeBoard(&game.BoardData)
	return &game
}
