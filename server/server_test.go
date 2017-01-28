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

// Create a new mux router.
var r = mux.NewRouter()

func init() {
	// Create a new chess server.
	s := New()

	// Add routes to the router.
	r.HandleFunc("/new", s.NewGameHandler)
	r.HandleFunc("/play/{id}/{auth}", s.PlayGameHandler)
}

func TestNewGameHandler(t *testing.T) {
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
	server := httptest.NewServer(r)
	defer server.Close()

	// Create a new game.
	resp, err := http.Get(server.URL + "/new")
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()

	path := fmt.Sprint("/play/invalid_game_id/invalid_auth")

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
	server := httptest.NewServer(r)
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
	path := fmt.Sprintf("/play/%s/invalid_auth_token", parts[0])

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
	server := httptest.NewServer(r)
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

	path := fmt.Sprintf("/play/%s/%s", gameId, whiteAuth)

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
}
