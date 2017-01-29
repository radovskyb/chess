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

	scanner := bufio.NewScanner(os.Stdin)
outer:
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		switch text {
		case "u":
			if err := b.UndoMove(); err != nil {
				fmt.Println(err)
				continue
			}
			b.Print()
			if hasCheck, color := b.HasCheck(); hasCheck {
				if b.InCheckmate(color) {
					fmt.Printf("%s is in checkmate\n", color)
					break
				}
				fmt.Printf("%s is in check\n", color)
			}
			if b.HasStalemate(b.Turn()) {
				fmt.Println("stalemate")
				break
			}
			continue
		case "p":
			history := b.History()
			if history != "" {
				fmt.Println(history)
			}
			continue
		case "q":
		inner:
			for {
				fmt.Print("are you sure you want to quit? (y/n): ")
				if !scanner.Scan() {
					break inner
				}
				text := strings.ToLower(strings.TrimSpace(scanner.Text()))
				switch {
				case strings.HasPrefix(text, "y"):
					return
				case strings.HasPrefix(text, "n"):
					break inner
				}
			}
			b.Print()
			continue
		}
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
		if err := b.MoveByLocation(loc1, loc2); err != nil {
			if err == engine.ErrInvalidLocation {
				fmt.Println("allowed formats: l1 l2 or l1l2 (e.g. a2 a4 or a2a4)")
			} else {
				fmt.Println(err)
			}
			continue
		}
		b.Print()
		if hasCheck, color := b.HasCheck(); hasCheck {
			if b.InCheckmate(color) {
				fmt.Printf("%s is in checkmate\n", color)
				break
			}
			fmt.Printf("%s is in check\n", color)
		}
		if b.HasStalemate(b.Turn()) {
			fmt.Println("stalemate")
			break
		}
		if mustPromote, _ := b.MustPromote(); mustPromote {
			fmt.Print("Promote pawn to? (k, r, b, q): ")
			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())
				var err error
				switch text {
				case "k":
					err = b.PromotePawn(engine.Knight)
				case "r":
					err = b.PromotePawn(engine.Rook)
				case "b":
					err = b.PromotePawn(engine.Bishop)
				case "q":
					err = b.PromotePawn(engine.Queen)
				default:
					fmt.Println("invalid piece, please choose between: k, r, b, q")
					continue
				}
				if err != nil {
					fmt.Println(err)
					continue
				}
				b.Print()
				if hasCheck, color := b.HasCheck(); hasCheck {
					if b.InCheckmate(color) {
						fmt.Printf("%s is in checkmate\n", color)
						break outer
					}
					fmt.Printf("%s is in check\n", color)
				}
				if b.HasStalemate(b.Turn()) {
					fmt.Println("stalemate")
					break outer
				}
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
