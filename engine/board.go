package engine

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidLocation   = errors.New("error: location string is invalid")
	ErrNoPieceAtPosition = errors.New("error: no piece at specified position")
	ErrOpponentsPiece    = errors.New("error: piece belongs to opponent")
	ErrInvalidPieceMove  = errors.New("error: invalid move for piece")
	ErrOccupiedPosition  = errors.New("error: position is already occupied")
	ErrMoveBlocked       = errors.New("error: move is blocked by another piece")
)

type Color uint8

const (
	Black Color = iota
	White
)

type Board struct {
	Turn       Color
	posToPiece map[Pos]*Piece
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
	return &Board{Turn: White, posToPiece: posToPiece}
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

func (b *Board) MoveByLocation(loc1, loc2 string) error {
	pos1, err := locToPos(loc1)
	if err != nil {
		return err
	}
	pos2, err := locToPos(loc2)
	if err != nil {
		return err
	}
	return b.MakeMove(pos1, pos2)
}

func (b *Board) move(piece *Piece, p1, p2 Pos) {
	// Move the piece to the new position.
	b.posToPiece[p2] = piece

	// Remove the piece from the old position.
	delete(b.posToPiece, p1)
}

func (b *Board) MakeMove(p1, p2 Pos) error {
	// Get the piece at position p1.
	piece, found := b.posToPiece[p1]
	if !found {
		return ErrNoPieceAtPosition
	}

	// Check that it's that piece's color's turn.
	if piece.Color != b.Turn {
		return ErrOpponentsPiece
	}

	// Check if the move is legal to make.
	if err := b.moveLegal(piece, p1, p2); err != nil {
		return err
	}

	// Move the piece on the board.
	b.move(piece, p1, p2)

	// Update who's turn it is.
	b.Turn ^= 1

	return nil
}

// moveLegal checks to see whether the specified move is legal to make or not.
func (b *Board) moveLegal(piece *Piece, p1, p2 Pos) error {
	// Get a list of all possible move positions that the
	// piece can move to without restrictions.
	//
	// TODO: Work out to check diagonal squares for pawns to take pieces.
	positions := getMovePositions(piece, p1)
	if _, ok := positions[p2]; !ok {
		return ErrInvalidPieceMove
	}

	// Check if there's a piece at position p2.
	//
	// TODO: Check King check line of sight.
	piece2, found := b.posToPiece[p2]

	// Check if piece2 is on the same team as piece.
	if found && piece.Color == piece2.Color {
		return ErrOccupiedPosition
	}

	// Check if the move from p1 to p2 is blocked by any other pieces.
	if b.moveBlocked(piece, p1, p2) {
		return ErrMoveBlocked
	}

	return nil
}

func (b *Board) moveBlocked(piece *Piece, p1, p2 Pos) bool {
	yd := 1
	if p1.Y > p2.Y {
		yd = -1
	}
	xd := 1
	if p1.X > p2.X {
		xd = -1
	}
	switch piece.Name {
	case Pawn:
		d := 1
		if piece.Color == Black {
			d = -1
		}
		if _, blocked := b.posToPiece[Pos{p1.X, p1.Y + (1 * d)}]; blocked {
			return true
		}
	case Rook:
		switch {
		case p1.Y != p2.Y:
			for y := p1.Y + (1 * yd); y != p2.Y; y = y + (1 * yd) {
				_, blocked := b.posToPiece[Pos{p2.X, y}]
				if blocked {
					return true
				}
			}
		case p1.X != p2.X:
			for x := p1.X + (1 * xd); x != p2.X; x = x + (1 * xd) {
				_, blocked := b.posToPiece[Pos{x, p2.Y}]
				if blocked {
					return true
				}
			}
		}
	}
	return false
}
