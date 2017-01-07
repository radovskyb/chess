package main

import (
	"fmt"
	"log"
	"time"
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
	time.Sleep(time.Second)

	if err := b.Move(Pos{1, 1}, Pos{1, 2}); err != nil {
		log.Fatalln(err)
	}
	b.Print()
	time.Sleep(time.Second)

	if err := b.Move(Pos{1, 2}, Pos{1, 3}); err != nil {
		log.Fatalln(err)
	}
	b.Print()
	time.Sleep(time.Second)

	if err := b.Move(Pos{2, 6}, Pos{2, 4}); err != nil {
		log.Fatalln(err)
	}
	b.Print()
	time.Sleep(time.Second)

	if err := b.Move(Pos{2, 4}, Pos{1, 3}); err != nil {
		log.Fatalln(err)
	}
	b.Print()
}
