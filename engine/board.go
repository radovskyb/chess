package engine

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

var (
	ErrInvalidLocation   = errors.New("error: location string is invalid")
	ErrNoPieceAtPosition = errors.New("error: no piece at specified position")
	ErrOpponentsPiece    = errors.New("error: piece belongs to opponent")
	ErrInvalidPieceMove  = errors.New("error: invalid move for piece")
	ErrOccupiedPosition  = errors.New("error: position is already occupied")
	ErrMoveBlocked       = errors.New("error: move is blocked by another piece")
	ErrMovingIntoCheck   = errors.New("error: move puts king in check")
	ErrMoveWhileInCheck  = errors.New("error: can't move piece while in check")
)

type Color uint8

const (
	Black Color = iota
	White
)

func (c Color) String() string {
	switch c {
	case Black:
		return "black"
	case White:
		return "white"
	default:
		return "invalid color"
	}
}

// A Board describes a chess board.
type Board struct {
	Turn       Color
	posToPiece map[Pos]*Piece
	kings      [2]Pos
	check      [2]bool
}

// HasCheck reports whether there is currently a king in check
// on the board.
func (b *Board) HasCheck() (bool, Color) {
	if b.check[Black] {
		return true, Black
	}
	if b.check[White] {
		return true, White
	}
	return false, 2 // 2 for invalid color.
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
	return &Board{
		Turn:       White,
		posToPiece: posToPiece,
		kings:      [2]Pos{White: {4, 0}, Black: {4, 7}},
	}
}

// clear removes all pieces from a board and is useful for testing.
func (b *Board) clear() {
	for k := range b.posToPiece {
		delete(b.posToPiece, k)
	}
}

var (
	printBlueBg = color.New(color.BgBlue).PrintfFunc()
	printCyanBg = color.New(color.BgCyan).PrintfFunc()
)

func (b *Board) Print() {
	fmt.Print("\033[H\033[2J\n")
	for i1 := 0; i1 < 8; i1++ {
		fmt.Print(color.RedString(" %d ", 8-i1))
		for i2 := 0; i2 < 8; i2++ {
			if piece, found := b.posToPiece[Pos{i2, 7 - i1}]; found {
				if i2%2 == i1%2 {
					printCyanBg("%s", piece)
				} else {
					printBlueBg("%s", piece)
				}
				continue
			}
			if i2%2 == i1%2 {
				printCyanBg("%3s", "")
			} else {
				printBlueBg("%3s", "")
			}
		}
		fmt.Println()
	}
	color.Red("%26s", "A  B  C  D  E  F  G  H")
}

// GetPieceAt returns the piece located at the string location.
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
