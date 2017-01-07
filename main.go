package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	Black = "b"
	White = "w"

	Pawn   = "P"
	Knight = "N"
	Bishop = "B"
	Rook   = "R"
	Queen  = "Q"
	King   = "K"
)

type Position struct {
	X, Y int
}

type Piece struct {
	Name  string
	Color string
}

func (p *Piece) String() string {
	return fmt.Sprintf("%2s%s", p.Color, p.Name)
}

type Board struct {
	Pieces map[Position]*Piece
}

func NewBoard() *Board {
	pieces := map[Position]*Piece{
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
		pieces[Position{i, 1}] = &Piece{Pawn, White}
		pieces[Position{i, 6}] = &Piece{Pawn, Black}
	}
	return &Board{Pieces: pieces}
}

func (b *Board) Print(w io.Writer) {
	for i1 := 0; i1 < 8; i1++ {
		for i2 := 0; i2 < 8; i2++ {
			if piece, found := b.Pieces[Position{i2, 7 - i1}]; found {
				fmt.Fprintf(w, "%s", piece)
				continue
			}
			fmt.Fprintf(w, "%3s", "--")
		}
		fmt.Println()
	}
}

// later return effects of move.
func (b *Board) Move(p1, p2 Position) error {
	piece, found := b.Pieces[p1]
	if !found {
		return fmt.Errorf("no piece at position %v", p1)
	}
	b.Pieces[p2] = piece
	delete(b.Pieces, p1)
	return nil
}

func main() {
	b := NewBoard()
	b.Print(os.Stdout)

	if err := b.Move(
		Position{1, 1},
		Position{1, 2},
	); err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	b.Print(os.Stdout)

	if err := b.Move(
		Position{1, 2},
		Position{1, 3},
	); err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	b.Print(os.Stdout)

	if err := b.Move(
		Position{2, 6},
		Position{2, 5},
	); err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	b.Print(os.Stdout)
}
