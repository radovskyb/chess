package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// Hard code the server's address for now.
const srvAddr = "http://localhost:9000"

// Create a new game and return the game's id, the current player's
// color with their auth string or any possible errors that occurred.
func NewGame() (string, string, string, error) {
	resp, err := http.Get(srvAddr + "/new")
	if err != nil {
		log.Fatalln(err)
	}
	slurp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", "", "", err
	}
	parts := strings.Split(string(slurp), ":")
	return parts[0], parts[1], parts[2], nil
}

// ConnectToGame connects to specified game for the player for color
// on the server and returns a websocket connection to the game.
func ConnectToGame(gameId, color, auth string) (*websocket.Conn, error) {
	path := fmt.Sprintf("/play/%s/%s/%s", gameId, color, auth)
	u := url.URL{
		Scheme: "ws",
		Host:   strings.TrimLeft(srvAddr, "http://"),
		Path:   path,
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	// Create a new game.
	gameId, color, auth, err := NewGame()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(gameId, color, auth)

	conn, err := ConnectToGame(gameId, color, auth)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println(string(msg))
	}
}
