package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/radovskyb/chess/engine"
)

func main() {
	b := engine.NewBoard()
	b.Print()

	// engine.FourMoveCheckmate(b) // play a four move checkmate animation.

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		locations := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(locations) != 2 {
			fmt.Println("enter format, loc1 loc2 (e.g. a2 a4)")
			continue
		}
		if err := b.MoveByLocation(locations[0], locations[1]); err != nil {
			if err == engine.ErrOpponentsPiece {
				fmt.Println("it's not your turn.")
			} else {
				fmt.Println(err)
			}
			continue
		}
		b.Print()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
