package main

import "fmt"

const (
	Black = "Black"
	White = "White"

	Pawn   = "Pawn"
	Knight = "Knight"
	Bishop = "Bishop"
	Rook   = "Rook"
	Queen  = "Queen"
	King   = "King"
)

type Position struct {
	X, Y int
}

type Piece struct {
	Name  string
	Color string
}

type Board struct {
	Pieces map[Position]Piece
}

func NewBoard() *Board {
	pieces := map[Position]Piece{
		{0, 0}: Piece{Rook, White},
		{0, 1}: Piece{Knight, White},
		{0, 2}: Piece{Bishop, White},
	}
	for i := 0; i < 8; i++ {
		pieces[Position{1, i}] = Piece{Pawn, White}
		pieces[Position{6, i}] = Piece{Pawn, White}
	}
	return &Board{Pieces: pieces}
}

func main() {
	b := NewBoard()
	piece, ok := b.Pieces[Position{6, 1}]
	if ok {
		fmt.Println(piece)
	} else {
		fmt.Println("no piece found at position 1, 1")
	}
}
