package main

import (
	"fmt"
	"log"
)

type Color uint8

const (
	Black Color = 1 << iota
	White

	Pawn   = "P"
	Knight = "N"
	Bishop = "B"
	Rook   = "R"
	Queen  = "Q"
	King   = "K"
)

type Pos struct {
	X, Y int
}

type Piece struct {
	Name string
	Color
}

func (p *Piece) AvailablePositions(cur Pos) map[Pos]struct{} {
	avail := make(map[Pos]struct{})
	switch {
	case p.Name == Pawn:
		avail[Pos{cur.X, cur.Y + 1}] = struct{}{}
	case p.Name == Knight:
		avail[Pos{cur.X + 2, cur.Y + 1}] = struct{}{}
		avail[Pos{cur.X - 2, cur.Y + 1}] = struct{}{}
		avail[Pos{cur.X + 2, cur.Y - 1}] = struct{}{}
		avail[Pos{cur.X - 2, cur.Y - 1}] = struct{}{}
		avail[Pos{cur.X + 1, cur.Y + 2}] = struct{}{}
		avail[Pos{cur.X - 1, cur.Y + 2}] = struct{}{}
		avail[Pos{cur.X + 1, cur.Y - 2}] = struct{}{}
		avail[Pos{cur.X - 1, cur.Y - 2}] = struct{}{}
	case p.Name == Bishop || p.Name == Queen: // fallthrough is not correct.
		for x, y := cur.X, cur.Y; x < 8 && y < 8; x, y = x+1, y+1 {
			avail[Pos{x, y}] = struct{}{}
		}
		for x, y := cur.X, cur.Y; x >= 0 && y < 8; x, y = x-1, y+1 {
			avail[Pos{x, y}] = struct{}{}
		}
		for x, y := cur.X, cur.Y; x < 8 && y >= 0; x, y = x+1, y-1 {
			avail[Pos{x, y}] = struct{}{}
		}
		for x, y := cur.X, cur.Y; x >= 0 && y >= 0; x, y = x-1, y-1 {
			avail[Pos{x, y}] = struct{}{}
		}
		fallthrough
	case p.Name == Rook || p.Name == Queen: // fallthrough is not correct.
		for x := cur.X; x < 8; x++ {
			avail[Pos{x, cur.Y}] = struct{}{}
		}
		for x := cur.X; x >= 0; x-- {
			avail[Pos{x, cur.Y}] = struct{}{}
		}
		for y := cur.Y; y < 8; y++ {
			avail[Pos{cur.X, y}] = struct{}{}
		}
		for y := cur.Y; y >= 0; y-- {
			avail[Pos{cur.X, y}] = struct{}{}
		}
	case p.Name == King:
		avail[Pos{cur.X + 1, cur.Y}] = struct{}{}
		avail[Pos{cur.X - 1, cur.Y}] = struct{}{}
		avail[Pos{cur.X, cur.Y + 1}] = struct{}{}
		avail[Pos{cur.X, cur.Y - 1}] = struct{}{}
	}
	return avail
}

func (p *Piece) String() string {
	color := "w"
	if p.Color == Black {
		color = "b"
	}
	return fmt.Sprintf("%2s%s", color, p.Name)
}

type Board struct {
	Turn   Color // who's turn it is.
	Pieces map[Pos]*Piece
	// pieceToPos map[*Piece]Pos
	// King Line Of Sights
}

func NewBoard() *Board {
	pieces := map[Pos]*Piece{
		// White
		{0, 0}: {Rook, White},
		{1, 0}: {Knight, White},
		{2, 0}: {Bishop, White},
		{3, 0}: {Queen, White},
		{4, 0}: {King, White},
		{5, 0}: {Bishop, White},
		{6, 0}: {Knight, White},
		{7, 0}: {Rook, White},
		// Black
		{0, 7}: {Rook, Black},
		{1, 7}: {Knight, Black},
		{2, 7}: {Bishop, Black},
		{3, 7}: {Queen, Black},
		{4, 7}: {King, Black},
		{5, 7}: {Bishop, Black},
		{6, 7}: {Knight, Black},
		{7, 7}: {Rook, Black},
	}
	for i := 0; i < 8; i++ {
		pieces[Pos{i, 1}] = &Piece{Pawn, White}
		pieces[Pos{i, 6}] = &Piece{Pawn, Black}
	}
	return &Board{Pieces: pieces}
}

func (b *Board) Print() {
	fmt.Print("\033[H\033[2J\n")
	for i1 := 0; i1 < 8; i1++ {
		for i2 := 0; i2 < 8; i2++ {
			if piece, found := b.Pieces[Pos{i2, 7 - i1}]; found {
				fmt.Printf("%s", piece)
				continue
			}
			fmt.Printf("%3s", "--")
		}
		fmt.Println()
	}
}

func (b *Board) Move(p1, p2 Pos) error {
	piece, found := b.Pieces[p1]
	if !found {
		return fmt.Errorf("no piece found at position %v", p1)
	}
	b.Pieces[p2] = piece // add to new pos
	delete(b.Pieces, p1) // delete piece at old pos
	return nil
}

func main() {
	b := NewBoard()
	b.Print()

	piece, found := b.Pieces[Pos{2, 0}]
	if !found {
		log.Fatalln("no bishop found")
	}
	// b.Move[Pos{2, 7}]
	fmt.Println(piece.AvailablePositions(Pos{2, 0}))
}

// func FourMoveCheckmate(b *Board) {
// 	// White 1
// 	if err := b.Move(Pos{4, 1}, Pos{4, 3}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// 	// Black 1
// 	if err := b.Move(Pos{4, 6}, Pos{4, 4}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// 	// White 2
// 	if err := b.Move(Pos{5, 0}, Pos{2, 3}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// 	// Black 2
// 	if err := b.Move(Pos{1, 7}, Pos{2, 5}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// 	// White 3
// 	if err := b.Move(Pos{3, 0}, Pos{7, 4}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// 	// Black 3
// 	if err := b.Move(Pos{3, 6}, Pos{3, 5}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// 	// White 4
// 	if err := b.Move(Pos{7, 4}, Pos{5, 6}); err != nil {
// 		log.Fatalln(err)
// 	}
// 	time.Sleep(time.Second)
// 	b.Print()
// }
