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
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "u" {
			if err := b.UndoMove(); err != nil {
				fmt.Println(err)
				continue
			}
			b.Print()
			continue
		}
		if text == "p" {
			history := b.History()
			if history != "" {
				fmt.Println(history)
			}
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
			if b.InCheckmate(color) {
				fmt.Printf("%s is in checkmate\n", color)
				break
			}
			fmt.Printf("%s is in check\n", color)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
