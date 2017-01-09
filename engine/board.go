package engine

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidLocation   = errors.New("error: location string is invalid")
	ErrNoPieceAtPosition = errors.New("error: no piece at specified position")
)

type Color uint8

const (
	Black Color = iota
	White
)

type Board struct {
	Turn       Color
	posToPiece map[Pos]*Piece
	// pieceToPos map[*Piece]Pos
}

func NewBoard() *Board {
	posToPiece := map[Pos]*Piece{
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
		posToPiece[Pos{i, 1}] = &Piece{Pawn, White}
		posToPiece[Pos{i, 6}] = &Piece{Pawn, Black}
	}
	// pieceToPos := make(map[*Piece]Pos)
	// for k, v := range posToPiece {
	// 	pieceToPos[v] = k
	// }
	// return &Board{posToPiece: posToPiece, pieceToPos: pieceToPos}
	return &Board{posToPiece: posToPiece}
}

func (b *Board) Print() {
	fmt.Print("\033[H\033[2J\n")
	for i1 := 0; i1 < 8; i1++ {
		for i2 := 0; i2 < 8; i2++ {
			if piece, found := b.posToPiece[Pos{i2, 7 - i1}]; found {
				fmt.Printf("%s", piece)
				continue
			}
			fmt.Printf("%3s", "--")
		}
		fmt.Println()
	}
}

// GetPieceAt either returns a piece located at the string location.
//
// If there's no piece at the specified location, or the location is invalid,
// an error is returned.
func (b *Board) GetPieceAt(loc string) (*Piece, error) {
	pos, err := locToPos(loc)
	if err != nil {
		return nil, err
	}
	piece, found := b.posToPiece[pos]
	if !found {
		return nil, ErrNoPieceAtPosition
	}
	return piece, nil
}

func (b *Board) Move(p1, p2 Pos) error {
	piece, found := b.posToPiece[p1]
	if !found {
		return ErrNoPieceAtPosition
	}
	b.posToPiece[p2] = piece // add to new pos
	// b.pieceToPos[piece] = p2 // update pieceToPos
	delete(b.posToPiece, p1) // delete piece at old pos
	return nil
}
