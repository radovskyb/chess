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

	// blockageMovesStr1 := "d2d4,e7e5,e1d2,f8b4"
	// blockageMovesStr2 := "a2a3,e7e5,a3a4,f8b4"
	// blockageMovesStr3 := "d2d4,e7e5,d4e5,d8e7,e5e6,e7e6a"
	// blockageMovesStr4 := "d2d4,e7e5,d4e5,d8e7,e5e6,e7e6,e2e4,e6e4,f1e2,f8b4"
	// setupForCheckStr1 := "e2e3,f7f5,a2a4,e8f7,d1e2,a7a5"
	setupForCheckStr2 := "e2e3,e7e5,b2b4,f8b4,a2a3,b4d2"
	// setupBoardStr1 := "g2g4,a7a5,g4g5,a5a4,g5g6,a4a3"
	// setupBoardStr2 := "a2a4,a7a5,a1a3,a8a6,a3c3,c7c5"
	// setupBoardStr3 := "a2a4,d7d5,a1a3,d5d4,a3g3,d4d3"

	moves := strings.Split(setupForCheckStr2, ",")
	for _, move := range moves {
		if err := b.MoveByLocation(move[0:2], move[2:]); err != nil {
			log.Fatalln(err)
		}
	}
	b.Print()

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
		if hasCheck, color := b.HasCheck(); hasCheck {
			fmt.Printf("%s is in check\n", color)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
