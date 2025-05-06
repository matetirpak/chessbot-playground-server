/*
API for existing games.
*/
package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/schema"

	data "github.com/matetirpak/chessbot-playground-server/internal/data"
	gl "github.com/matetirpak/chessbot-playground-server/internal/game_logic"
)

// GetGame godoc
//
//		@Summary			Extract board, possible moves or wait for a turn notification
//		@Description		Requesting the board state requires the 'moveidx' parameter. Requesting possible moves of a piece requires 'row' and 'col'. 'turn' holds the request till it's the players turn.
//		@Tags				game
//		@Accept				json
//		@Produce			json
//	    @Security 			BearerAuth
//		@Param				moveidx		query		int		false			"Index of move (for board state)"
//		@Param				boardid		query		int		true			"Board ID"
//		@Param				color		query		string	true			"Color ('w' or 'b')"
//		@Param				row			query		int		false			"Row of piece (for moves)"
//		@Param				col			query		int		false			"Column of piece (for moves)"
//		@Param				reqtype		query		string	true			"Request type: 'state', 'turn' or 'moves'"
//		@Success           	200      	{object}	interface{}  			"Returns one of: BoardState (reqtype=state), []RespMove (reqtype=moves), or {} (reqtype=turn)". Defined at internal/api/structs.go
//		@Failure			400			{string}	string					"Bad request (invalid parameters or out-of-range index)"
//		@Failure			401			{string}	string					"Unauthorized (missing/invalid token)"
//		@Failure			404			{string}	string					"Not found – Game does not exist"
//		@Failure			408			{string}	string					"Timeout waiting for turn"
//		@Failure			500			{string}	string					"Internal server error during move generation"
//		@Router				/chessserver/v1/game [get]
func GetGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get Bearer token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	var req ReqGetGame
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse query params: %v", err), http.StatusBadRequest)
		return
	}

	if req.ReqType != "state" && req.ReqType != "turn" && req.ReqType != "moves" {
		http.Error(w, "\"reqtype\" has to be \"state\", \"turn\" or \"moves\"", http.StatusBadRequest)
		return
	}

	data.GamesMapMu.RLock()
	game, exists := data.GamesMap[req.BoardID]
	nSteps := len(game.BoardData)
	data.GamesMapMu.RUnlock()
	if !exists {
		http.Error(w, fmt.Sprintf("Game with index %d doesn't exist.", req.BoardID), http.StatusNotFound)
		return
	}

	if req.ReqType == "turn" {
		// Check session access
		success := verifyGameAccess(w, req.BoardID, token)
		if !success {
			return
		}
	} else {
		// Check player related board access
		success := verifyBoardAccess(w, game, req.Color, token)
		if !success {
			return
		}
	}

	switch req.ReqType {
	case "state":
		idx := int(req.Moveidx)
		if idx == -1 {
			idx = nSteps - 1
		}
		if idx > nSteps-1 {
			http.Error(w, fmt.Sprintf("Board at index %d does not exist.", idx), http.StatusNotFound)
			return
		}
		game.Mu.RLock()
		resp := game.BoardData[idx]
		game.Mu.RUnlock()

		json.NewEncoder(w).Encode(resp)
	case "turn":
		w.Header().Set("Connection", "keep-alive")

		timeout := time.After(1 * time.Hour)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				http.Error(w, "Timeout waiting for turn.", http.StatusRequestTimeout)
				return
			case <-ticker.C:
				// Check for current player's turn
				game.Mu.RLock()
				currentTurn := game.BoardData[len(game.BoardData)-1].TurnColor
				game.Mu.RUnlock()

				if currentTurn == req.Color {
					w.WriteHeader(http.StatusOK)
					return
				}
			case <-r.Context().Done():
				http.Error(w, "Server exited while waiting for turn.", http.StatusRequestTimeout)
				return
			}
		}
	case "moves":
		var moves *gl.Move

		game.Mu.RLock()
		bstate := game.BoardData[len(game.BoardData)-1]
		game.Mu.RUnlock()

		gl.GenerateMovesForPiece(int(req.Row), int(req.Col), &bstate, &moves)
		validMoves, err := gl.FilterInvalidMoves(moves, &bstate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to validate generated moves with error: %v", err), http.StatusInternalServerError)
		}
		movesList := llToArray(validMoves)

		json.NewEncoder(w).Encode(movesList)
	default:
		http.Error(w, "Invalid parameter \"reqtype\"", http.StatusBadRequest)
	}
}

// PutGame godoc
//
//	@Summary		Applies an action to a game
//	@Description	Applies a specified (valid) move, random move, or forfeits the game. Applying a specific move requires the move variable, e.g. "e2 e4".
//	@Tags			game
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		ReqPutGame	true	"reqtype: 'move' (requires 'move' variable), 'randommove', 'forfeit'"
//	@Success		200		{string}	string				"Success (No Content)"
//	@Failure		400		{string}	string				"Bad request (invalid parameters, move, or game state)"
//	@Failure		401		{string}	string				"Unauthorized (missing or invalid token)"
//	@Failure		404		{string}	string				"Not found – Game does not exist"
//	@Failure		500		{string}	string				"Internal server error during move processing"
//	@Router			/chessserver/v1/game [put]
func PutGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get Bearer token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	var req ReqPutGame
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	if req.ReqType != "forfeit" && req.ReqType != "move" && req.ReqType != "randommove" {
		http.Error(w, "\"reqtype\" has to be \"forfeit\", \"move\" or \"randommove\"", http.StatusBadRequest)
		return
	}

	data.GamesMapMu.RLock()
	game, exists := data.GamesMap[req.BoardID]
	data.GamesMapMu.RUnlock()
	if !exists {
		http.Error(w, fmt.Sprintf("Game with index %d doesn't exist.", req.BoardID), http.StatusNotFound)
		return
	}

	success2 := verifyBoardAccess(w, game, req.Color, token)
	if !success2 {
		return
	}

	if req.ReqType == "forfeit" {
		game.Mu.Lock()
		game.BoardData[len(game.BoardData)-1].TurnColor = "n"
		if req.Color == "w" {
			game.Winner = "b"
			game.BoardData[len(game.BoardData)-1].Winner = "b"
		} else {
			game.Winner = "w"
			game.BoardData[len(game.BoardData)-1].Winner = "w"
		}
		game.Mu.Unlock()

		w.WriteHeader(http.StatusOK)
		return
	}

	game.Mu.RLock()
	gameStarted := game.Started
	gameWinner := game.Winner
	latestBoardState := game.BoardData[len(game.BoardData)-1]
	game.Mu.RUnlock()

	if !gameStarted {
		http.Error(w, "Can't apply move. Game has not started.", http.StatusBadRequest)
		return
	}

	if gameWinner != "n" {
		http.Error(w, "Can't apply move. Game has ended.", http.StatusBadRequest)
		return
	}

	if latestBoardState.TurnColor != req.Color {
		http.Error(w, "Can't apply move. It's not the players turn.", http.StatusBadRequest)
		return
	}

	var move gl.Move
	if req.ReqType == "move" {
		move, err = gl.StringToMoveStruct(req.Move, rune(req.Color[0]))
		if err != nil {
			http.Error(w, "Move format is invalid.", http.StatusBadRequest)
			return
		}

		// Check validity of move
		err = gl.ValidateMove(&move, &latestBoardState)
		if err != nil {
			http.Error(w, fmt.Sprintf("Move is invalid with error: %v", err), http.StatusBadRequest)
			return
		}
	}
	if req.ReqType == "randommove" {
		moves := gl.AllPossibleMoves(rune(req.Color[0]), &latestBoardState, nil)
		validMoves, err := gl.FilterInvalidMoves(moves, &latestBoardState)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to validate generated moves with error: %v", err), http.StatusInternalServerError)
			return
		}
		nMoves := gl.NMoves(validMoves)
		if nMoves == 0 {
			http.Error(w, "No moves found.", http.StatusInternalServerError)
			return
		}
		randomIdx := rand.Intn(nMoves)
		move = *gl.MoveAt(validMoves, randomIdx)
	}

	game.Mu.Lock()
	newBstate := gl.MakeMove(&move, game.BoardData[len(game.BoardData)-1], true)
	game.Winner = newBstate.Winner
	game.BoardData = append(game.BoardData, newBstate)
	game.Mu.Unlock()

	w.WriteHeader(http.StatusOK)
}
