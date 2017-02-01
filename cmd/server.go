package main

import (
	"log"
	"net/http"

	"github.com/radovskyb/chess/server"
)

func main() {
	r := server.GetRouter()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
