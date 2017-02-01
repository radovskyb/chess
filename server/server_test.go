package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func newRouter() *mux.Router {
	// Create a new mux router.
	r := mux.NewRouter()

	// Create a new chess server.
	s := New()

	// Add routes to the router.
	r.HandleFunc("/new", s.NewGameHandler)
	r.HandleFunc("/play/{id}/{color}/{auth}", s.PlayGameHandler)

	return r
}

func TestNewGameHandler(t *testing.T) {
	r := newRouter()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rr.Code)
	}

	parts := strings.Split(rr.Body.String(), ":")
	if len(parts) != 3 {
		t.Errorf("expected body to contain 3 parts, got %d parts", len(parts))
	}

	gameId1, color := parts[0], parts[1]
	if color != "white" {
		t.Errorf("expected color to be white, got %s", color)
	}

	// Reset rr's byte buffer's body.
	rr.Body.Reset()

	// Make the same request again.
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rr.Code)
	}

	parts = strings.Split(rr.Body.String(), ":")
	if len(parts) != 3 {
		t.Errorf("expected body to contain 3 parts, got %d parts", len(parts))
	}

	gameId2, color := parts[0], parts[1]
	if color != "black" {
		t.Errorf("expected color to be black, got %s", color)
	}

	if gameId1 != gameId2 {
		t.Error("expected gameId1 and gameId2 to match")
	}
}

func TestConnectWithWrongGameId(t *testing.T) {
	server := httptest.NewServer(newRouter())
	defer server.Close()

	// Create a new game.
	resp, err := http.Get(server.URL + "/new")
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()

	path := fmt.Sprint("/play/invalid_game_id/white/invalid_auth")

	u := url.URL{
		Scheme: "ws",
		Host:   strings.TrimLeft(server.URL, "http://"),
		Path:   path,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	mt, msg, err := c.ReadMessage()
	if err != nil {
		t.Error(err)
	}

	if mt != websocket.TextMessage {
		t.Errorf("expected mt to be a close message, got %v", mt)
	}

	if string(msg) != ErrGameIdNotFound.Error() {
		t.Errorf("expected msg to be a ErrGameIdNotFound error, got %s", msg)
	}
}

func TestConnectWithWrongAuthToken(t *testing.T) {
	server := httptest.NewServer(newRouter())
	defer server.Close()

	// Create a new game.
	resp, err := http.Get(server.URL + "/new")
	if err != nil {
		t.Error(err)
	}
	slurp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
	}

	parts := strings.Split(string(slurp), ":")
	path := fmt.Sprintf("/play/%s/white/invalid_auth_token", parts[0])

	u := url.URL{
		Scheme: "ws",
		Host:   strings.TrimLeft(server.URL, "http://"),
		Path:   path,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	mt, msg, err := c.ReadMessage()
	if err != nil {
		t.Error(err)
	}

	if mt != websocket.TextMessage {
		t.Errorf("expected mt to be a close message, got %v", mt)
	}

	if string(msg) != ErrInvalidAuthToken.Error() {
		t.Errorf("expected msg to be a ErrInvalidAuthToken error, got %s", msg)
	}
}

func TestConnectWithCorrectInfoSuccessful(t *testing.T) {
	server := httptest.NewServer(newRouter())
	defer server.Close()

	// Create a new game.
	resp, err := http.Get(server.URL + "/new")
	if err != nil {
		t.Error(err)
	}
	slurp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
	}

	parts := strings.Split(string(slurp), ":")
	gameId, whiteAuth := parts[0], parts[2]

	path := fmt.Sprintf("/play/%s/white/%s", gameId, whiteAuth)

	u := url.URL{
		Scheme: "ws",
		Host:   strings.TrimLeft(server.URL, "http://"),
		Path:   path,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error(err)
	}
	c.Close()
}

func connectToGame(server *httptest.Server, t *testing.T) (string, string, *websocket.Conn) {
	// Create a new game.
	resp, err := http.Get(server.URL + "/new")
	if err != nil {
		t.Error(err)
	}
	slurp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
	}

	parts := strings.Split(string(slurp), ":")
	gameId, color, auth := parts[0], parts[1], parts[2]

	path := fmt.Sprintf("/play/%s/%s/%s", gameId, color, auth)

	u := url.URL{
		Scheme: "ws",
		Host:   strings.TrimLeft(server.URL, "http://"),
		Path:   path,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error(err)
	}

	return gameId, auth, c
}

func TestWaitForSecondPlayerToJoin(t *testing.T) {
	server := httptest.NewServer(newRouter())
	defer server.Close()

	_, _, whiteConn := connectToGame(server, t)
	defer whiteConn.Close()

	mt, msg, err := whiteConn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if mt == websocket.CloseMessage {
		t.Error("expected text message from websocket, got close message")
	}
	if string(msg) != "waiting for player to join" {
		t.Errorf("expected waiting for player message, got %s", msg)
	}

	// Connect black to the game.
	gameId, blackAuth, blackConn := connectToGame(server, t)

	mt, msg, err = whiteConn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if mt == websocket.CloseMessage {
		t.Error("expected text message from websocket, got close message")
	}
	if string(msg) != "starting game" {
		t.Errorf("expected starting game message, got %s", msg)
	}

	blackConn.Close()

	mt, msg, err = whiteConn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if mt == websocket.CloseMessage {
		t.Error("expected text message from websocket, got close message")
	}
	if string(msg) != "waiting for player to join" {
		t.Errorf("expected waiting for player message, got %s", msg)
	}

	// Re-Connect black to the game.
	path := fmt.Sprintf("/play/%s/black/%s", gameId, blackAuth)

	u := url.URL{
		Scheme: "ws",
		Host:   strings.TrimLeft(server.URL, "http://"),
		Path:   path,
	}

	blackConn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error(err)
	}
	blackConn.Close()

	mt, msg, err = whiteConn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if mt == websocket.CloseMessage {
		t.Error("expected text message from websocket, got close message")
	}
	if string(msg) != "starting game" {
		t.Errorf("expected starting game message, got %s", msg)
	}
}
