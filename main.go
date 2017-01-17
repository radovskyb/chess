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
	// if err := setupBoardForCheck(b); err != nil {
	// 	log.Fatalln(err)
	// }
	setupBoardForCheck(b)
	// setupBoardForBlockage(b)
	// setupBoardForBlockage2(b)
	// setupBoardForBlockage3(b)
	// setupBoardForBlockage4(b)

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

func setupBoardForBlockage4(b *engine.Board) error {
	b.MoveByLocation("d2", "d4")
	b.MoveByLocation("e7", "e5")
	b.MoveByLocation("e1", "d2")
	b.MoveByLocation("f8", "b4")
	// b.MoveByLocation("d8", "e7")
	// b.MoveByLocation("e5", "e6")
	// b.MoveByLocation("e7", "e6")
	// b.MoveByLocation("e2", "e4")
	// b.MoveByLocation("e6", "e4")
	// b.MoveByLocation("f1", "e2")
	b.Print()
	return nil
}

func setupBoardForBlockage3(b *engine.Board) error {
	b.MoveByLocation("d2", "d4")
	b.MoveByLocation("e7", "e5")
	b.MoveByLocation("d4", "e5")
	b.MoveByLocation("d8", "e7")
	b.MoveByLocation("e5", "e6")
	b.MoveByLocation("e7", "e6")
	b.MoveByLocation("e2", "e4")
	b.MoveByLocation("e6", "e4")
	b.MoveByLocation("f1", "e2")
	b.MoveByLocation("f8", "b4")
	// b.MoveByLocation("d1", "d2")
	// b.MoveByLocation("e4", "e2")
	b.Print()
	return nil
}

func setupBoardForBlockage2(b *engine.Board) error {
	b.MoveByLocation("d2", "d4")
	b.MoveByLocation("e7", "e5")
	b.MoveByLocation("d4", "e5")
	b.MoveByLocation("d8", "e7")
	b.MoveByLocation("e5", "e6")
	b.MoveByLocation("e7", "e6")
	// b.MoveByLocation("e2", "e3") // <-- causing unwanted check.
	b.Print()
	return nil
}

func setupBoardForBlockage(b *engine.Board) error {
	b.MoveByLocation("a2", "a3")
	b.MoveByLocation("e7", "e5")
	b.MoveByLocation("a3", "a4")
	b.MoveByLocation("f8", "b4")
	b.Print()
	return nil
}

func setupBoardForCheck(b *engine.Board) error {
	if err := b.MoveByLocation("e2", "e3"); err != nil {
		return err
	}
	if err := b.MoveByLocation("f7", "f5"); err != nil {
		return err
	}
	if err := b.MoveByLocation("a2", "a4"); err != nil {
		return err
	}
	if err := b.MoveByLocation("e8", "f7"); err != nil {
		return err
	}
	if err := b.MoveByLocation("d1", "e2"); err != nil {
		return err
	}
	if err := b.MoveByLocation("a7", "a5"); err != nil {
		return err
	}
	b.Print()
	return nil
}

func setupBoard1(b *engine.Board) error {
	b.MoveByLocation("g2", "g4")
	b.MoveByLocation("a7", "a5")
	b.MoveByLocation("g4", "g5")
	b.MoveByLocation("a5", "a4")
	b.MoveByLocation("g5", "g6")
	b.MoveByLocation("a4", "a3")
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

func setupBoard3(b *engine.Board) error {
	b.MoveByLocation("a2", "a4")
	b.MoveByLocation("d7", "d5")
	b.MoveByLocation("a1", "a3")
	b.MoveByLocation("d5", "d4")
	b.MoveByLocation("a3", "g3")
	b.MoveByLocation("d4", "d3")
	b.Print()
	return nil
}
