package server

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/radovskyb/chess/engine"
)

var (
	ErrGameIdNotFound   = errors.New("game id not found")
	ErrInvalidAuthToken = errors.New("invalid auth token")
)

var upgrader = websocket.Upgrader{}

func init() {
	// Create a new mux router.
	r := mux.NewRouter()

	s := New()

	r.HandleFunc("/new", s.NewGameHandler)
	r.HandleFunc("/play/{id}/{auth}", s.PlayGameHandler)

	http.Handle("/", r)
}

// A game holds a chess board and authentication tokens
// for both players.
type game struct {
	// colors holds authentication tokens for the player
	// for white and the player for black.
	colors [2]string

	// moves is a channel that always holds the info for the last
	// move that was made and only gets cleared in between when the
	// move is updated on the opponent's client's game board and an
	// ack has been sent from that client back to the server.
	moves chan *engine.MoveInfo

	// board holds the chess board.
	board *engine.Board

	// turn is a channel used for polling so client's can receive a
	// signal when the player has finished making their move.
	turn chan engine.Color
}

type Server struct {
	mu *sync.Mutex // mu protects the following.

	// games holds a map of game id's to games.
	games map[string]*game

	// unmatched contains a list of games currently containing only
	// a single player and the time that the game was created.
	unmatched *list.List
	// umLookup is a reverse lookup for the unmatched game list elements
	// based on the game's id.
	umLookup map[string]*list.Element
}

// unmatchedElem holds a game id and the time the game
// was created and put into the unmatched list.
type unmatchedElem struct {
	gameId  string
	created time.Time
}

func New() *Server {
	return &Server{
		mu:        new(sync.Mutex),
		games:     make(map[string]*game),
		unmatched: list.New(),
		umLookup:  make(map[string]*list.Element),
	}
}

// PlayGameHandler connects to a game by it's id. It opens a websocket to the client.
func (s *Server) PlayGameHandler(w http.ResponseWriter, r *http.Request) {
	// Open WebSocket connection and verify authentication for client.
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	// ========================================

	// Validate that the game id exists and the client has
	// correct authentication to play the game.
	vars := mux.Vars(r)
	gameId, auth := vars["id"], vars["auth"]

	// Check that the game id exists.
	game, found := s.games[gameId]
	if !found {
		c.WriteMessage(websocket.TextMessage, []byte(ErrGameIdNotFound.Error()))
		return
	}

	// Check that the auth token exists for either the white or black player.
	if game.colors[engine.White] != auth ||
		game.colors[engine.Black] != auth {
		c.WriteMessage(websocket.TextMessage, []byte(ErrInvalidAuthToken.Error()))
		return
	}

	// ========================================

	// If there's only 1 player, wait for a second player to join.

	// ========================================

	// If it's the client's turn, tell the client.

	// ========================================

	// If a move is made, check if it's the right client's turn and if so,
	// make a move on the game's board and then place the move on the
	// game's move channel for the client to receive either right now, or
	// the next time the client reconnects to the game. If it's not the
	// sending client's turn, send them an error message that it's not their
	// turn.

	// ========================================

	// Once the move has been sent to the client, wait for an ack to be received
	// from the other client, which is to let the server know that their board has
	// been updated.

	// ========================================

	// When a client disconnects, if the other client is still connected, send a
	// signal on a channel so the other client can be notified, otherwise just
	// return and close the current connection.

	// ========================================

	// If a quit message is received, notify the other client and do a cleanup.
	// (CLEANUP): Take contents of QuitGameHandler.
}

// NewGameHandler creates a new game in the server's games cache and
// returns the game's id to the client.
func (s *Server) NewGameHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate an auth token for the player for the new game.
	auth, err := newToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there's already an unmatched game, match the current client with it.
	//
	// For now, unmatched games will only contain white players since new games
	// are automatically added with a white player first, so by default if the game
	// is unmatched, set the player to the black player.
	if s.unmatched.Len() > 0 {
		unmatched := s.unmatched.Front()

		gameId := unmatched.Value.(*unmatchedElem).gameId

		// Write the game id string, auth token and color to the client.
		//
		// The joining player will be black.
		if _, err := fmt.Fprintf(w, "%s:%s:%s",
			gameId,
			engine.Black,
			auth,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Remove the game from the unmatched game lookup map.
		delete(s.umLookup, gameId)

		// Remove the game from the unmatched list.
		s.unmatched.Remove(unmatched)

		// Set the auth for black in the game.
		s.games[gameId].colors[engine.Black] = auth

		return
	}

	// Create a new game since one doesn't already exist.

	// Generate a token for the game id.
	id, err := newToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add a new game to s.games with the current client being set
	// to play for white on the board.
	s.games[id] = &game{
		colors: [2]string{engine.White: auth},
		board:  engine.NewBoard(),
		moves:  make(chan *engine.MoveInfo, 1),
		turn:   make(chan engine.Color, 1),
	}

	// Add the game to the unmatched list.
	umGame := s.unmatched.PushBack(&unmatchedElem{
		gameId: id, created: time.Now(),
	})
	// Add the unmatched game to the unmatched game lookup map.
	s.umLookup[id] = umGame

	// Write the game id string, auth token and color to the client.
	if _, err := fmt.Fprintf(w, "%s:%s:%s",
		id,
		engine.White,
		auth,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// newToken returns a cryptographically secure and url safe
// encoded token.
func newToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Encode the random bytes to a url safe encoded string and return it.
	return base64.URLEncoding.EncodeToString(b), nil
}
