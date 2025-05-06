package api

// Create new game
type ReqPostSessions struct {
	Name string `json:"name"`
}
type RespPostSessions struct {
	BoardID  int32  `json:"boardid"`
	Password string `json:"password"`
}

// Delete an ongoing game
type ReqDeleteSessions struct {
	BoardID int32 `json:"boardid"`
}

// Get all sessions
type RespGetSessions struct {
	Games []GameNameAndID `json:"games"`
}
type GameNameAndID struct {
	Name    string `json:"name"`
	BoardID int32  `json:"boardid"`
}

// Entry a game
type ReqPutSessions struct {
	BoardID int32  `json:"boardid"`
	Color   string `json:"color"`
}
type RespPutSessions struct {
	Token string `json:"token"`
}

// Get game data
type ReqGetGame struct {
	Moveidx int32  `schema:"moveidx"`
	BoardID int32  `schema:"boardid"`
	Color   string `schema:"color"`
	Row     int32  `schema:"row"`
	Col     int32  `schema:"col"`
	ReqType string `schema:"reqtype"` // "state" "turn" "moves"
}

type RespMove struct {
	From    [2]int `json:"from"`
	To      [2]int `json:"to"`
	Capture bool   `json:"capture"`
}

// Apply move
type ReqPutGame struct {
	BoardID int32  `json:"boardid"`
	Color   string `json:"color"`
	Move    string `json:"move,omitempty"`
	ReqType string `json:"reqtype"` // "forfeit" "move" "randommove"
}
