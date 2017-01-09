package main

import (
	"fmt"
	"log"

	"github.com/radovskyb/chess/engine"
)

func main() {
	b := engine.NewBoard()
	b.Print()

	piece, err := b.GetPieceAt("a1")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\n%s\n", piece)

	// engine.FourMoveCheckmate(b)
}
