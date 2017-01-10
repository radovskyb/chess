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
	// if err := setupBoard1(b); err != nil {
	// 	log.Fatalln(err)
	// }
	// if err := setupBoard2(b); err != nil {
	// 	log.Fatalln(err)
	// }

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		var loc1, loc2 string
		switch len(text) {
		case 4:
			loc1, loc2 = text[0:2], text[2:4]
		case 5:
			locations := strings.Split(strings.TrimSpace(scanner.Text()), " ")
			if len(locations) == 2 {
				loc1, loc2 = locations[0], locations[1]
			}
		}
		if loc1 == "" || loc2 == "" {
			fmt.Println("allowed formats: l1 l2 or l1l2 (e.g. a2 a4 or a2a4)")
			continue
		}
		if err := b.MoveByLocation(loc1, loc2); err != nil {
			fmt.Println(err)
			continue
		}
		b.Print()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func setupBoard1(b *engine.Board) error {
	b.MoveByLocation("a2", "a4")
	b.MoveByLocation("a7", "a5")
	b.MoveByLocation("a1", "a3")
	b.MoveByLocation("a8", "a6")
	b.MoveByLocation("b2", "b3")
	b.MoveByLocation("b7", "b6")
	b.MoveByLocation("h2", "h4")
	b.MoveByLocation("h7", "h5")
	b.Print()
	return nil
}

func setupBoard2(b *engine.Board) error {
	b.MoveByLocation("a2", "a4")
	b.MoveByLocation("a7", "a5")
	b.MoveByLocation("a1", "a3")
	b.MoveByLocation("a8", "a6")
	b.MoveByLocation("a3", "c3")
	b.MoveByLocation("c7", "c5")
	b.Print()
	return nil
}
