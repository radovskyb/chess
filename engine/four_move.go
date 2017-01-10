package engine

import (
	"log"
	"time"
)

func FourMoveCheckmate(b *Board) {
	// White 1
	if err := b.MakeMove(Pos{4, 1}, Pos{4, 3}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
	// Black 1
	if err := b.MakeMove(Pos{4, 6}, Pos{4, 4}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
	// White 2
	if err := b.MakeMove(Pos{5, 0}, Pos{2, 3}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
	// Black 2
	if err := b.MakeMove(Pos{1, 7}, Pos{2, 5}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
	// White 3
	if err := b.MakeMove(Pos{3, 0}, Pos{7, 4}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
	// Black 3
	if err := b.MakeMove(Pos{3, 6}, Pos{3, 5}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
	// White 4
	if err := b.MakeMove(Pos{7, 4}, Pos{5, 6}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)
	b.Print()
}
