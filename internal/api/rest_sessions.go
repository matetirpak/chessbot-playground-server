/*
API for session management.
*/
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	data "github.com/matetirpak/chessbot-playground-server/internal/data"
)

// DeleteSessions godoc
//
//	@Summary		Deletes a session
//	@Description	Given the board-id, deletes the session if the correct session token was provided.
//	@Tags			sessions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		ReqDeleteSessions	true	"Request payload with board-id"
//	@Success		200		{string}	string						"Success (No Content)"
//	@Failure		400		{string}	string						"Bad request (invalid JSON body)"
//	@Failure		401		{string}	string						"Unauthorized (missing or invalid session token)"
//	@Failure		404		{string}	string						"Not found – Game session does not exist"
//	@Router			/chessserver/v1/sessions [delete]
func DeleteSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get Bearer token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	password := strings.TrimPrefix(authHeader, "Bearer ")

	var req ReqDeleteSessions
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Decoding failed.\n")
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	success := verifyGameAccess(w, req.BoardID, password)
	if !success {
		return
	}
	data.GamesMapMu.Lock()
	delete(data.GamesMap, req.BoardID)
	data.GamesMapMu.Unlock()

	w.WriteHeader(http.StatusOK)
}

// GetSessions godoc
//
//	@Summary		Displays all existing sessions
//	@Description	Returns a list of all existing sessions. No request body or parameters are required.
//	@Tags			sessions
//	@Produce		json
//	@Success 		200 	{object} 	RespGetSessions 			"List of all game sessions"
//	@Router			/chessserver/v1/sessions [get]
func GetSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var resp RespGetSessions

	// Iterate through the GamesMap and populate the response
	data.GamesMapMu.RLock()
	for _, game := range data.GamesMap {
		extracted_game := GameNameAndID{
			Name:    game.Name,
			BoardID: game.ID,
		}
		// Append the response to the slice
		resp.Games = append(resp.Games, extracted_game)
	}
	data.GamesMapMu.RUnlock()

	json.NewEncoder(w).Encode(resp)
}

// PostSessions godoc
//
//	@Summary		Creates a new session
//	@Description	Initializes a new session in the server. The response contains an ID and password.
//	@Tags			sessions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ReqPostSessions		true	"Request payload with desired session name"
//	@Success 		200 	{object} 	RespPostSessions 			"Session/Board ID and password"
//	@Failure		400		{string}	string						"Bad request (invalid JSON body)"
//	@Router			/chessserver/v1/sessions [post]
func PostSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var req ReqPostSessions
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	newGame := initializeNewGame(req.Name)

	data.GamesMapMu.Lock()
	data.GamesMap[newGame.ID] = newGame
	data.GamesMapMu.Unlock()

	resp := RespPostSessions{BoardID: newGame.ID, Password: newGame.Password}
	json.NewEncoder(w).Encode(resp)
}

// PutSessions godoc
//
//	@Summary		Register as a player in a session
//	@Description	Registers a player (white or black) to an existing game session using the board ID and a session password.
//	@Tags			sessions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		ReqPutSessions		true	"Request payload with session access data and desired color"
//	@Success 		200 	{object} 	RespPutSessions 			"Color-specific Player token"
//	@Failure		400		{string}	string				"Bad request – Invalid JSON or color value"
//	@Failure		401		{string}	string				"Unauthorized – Missing or invalid bearer token"
//	@Failure		403		{string}	string				"Forbidden – Game is full or color already taken"
//	@Failure		404		{string}	string				"Not found – Game session does not exist"
//	@Router			/chessserver/v1/sessions [put]
func PutSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var req ReqPutSessions
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Get Bearer token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	password := strings.TrimPrefix(authHeader, "Bearer ")

	success := verifyGameAccess(w, req.BoardID, password)
	if !success {
		return
	}

	data.GamesMapMu.RLock()
	var game *data.Game = data.GamesMap[req.BoardID]
	data.GamesMapMu.RUnlock()

	game.Mu.RLock()
	hasWPlayer := game.HasWPlayer
	hasBPlayer := game.HasBPlayer
	game.Mu.RUnlock()

	if hasWPlayer && hasBPlayer {
		http.Error(w, "Game is already full.", http.StatusForbidden)
		return
	}

	token := generateToken()
	resp := RespPutSessions{Token: token}

	switch req.Color {
	case "w":
		if hasWPlayer {
			http.Error(w, "White is already taken.", http.StatusForbidden)
			return
		}
		game.Mu.Lock()
		game.HasWPlayer = true
		game.WPlayerToken = token
		if game.HasBPlayer {
			game.Started = true
			game.BoardData[0].TurnColor = "w"
		}
		game.Mu.Unlock()

		json.NewEncoder(w).Encode(resp)

	case "b":
		if hasBPlayer {
			http.Error(w, "Black is already taken.", http.StatusForbidden)
			return
		}
		game.Mu.Lock()
		game.HasBPlayer = true
		game.BPlayerToken = token
		if game.HasWPlayer {
			game.Started = true
			game.BoardData[0].TurnColor = "w"
		}
		game.Mu.Unlock()

		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Invalid color. Enter 'w' or 'b'", http.StatusBadRequest)
		return
	}
}
