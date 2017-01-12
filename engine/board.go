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
)

type Color uint8

const (
	Black Color = iota
	White
)

type Board struct {
	Turn       Color
	posToPiece map[Pos]*Piece
	kings      [2]Pos
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
					printBlueBg("%s", piece)
				} else {
					printCyanBg("%s", piece)
				}
				continue
			}
			if i2%2 == i1%2 {
				printBlueBg("%3s", "")
			} else {
				printCyanBg("%3s", "")
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

func (b *Board) MoveByLocation(loc1, loc2 string) error {
	pos1, err := locToPos(loc1)
	if err != nil {
		return err
	}
	pos2, err := locToPos(loc2)
	if err != nil {
		return err
	}
	return b.Move(pos1, pos2)
}

// moveByLocation is a convenience function used for setting up
// boards for testing by moving pieces by locations and also
// avoiding any checks for move legality.
func (b *Board) moveByLocation(loc1, loc2 string) error {
	pos1, err := locToPos(loc1)
	if err != nil {
		return err
	}
	pos2, err := locToPos(loc2)
	if err != nil {
		return err
	}
	piece, found := b.posToPiece[pos1]
	if !found {
		return ErrNoPieceAtPosition
	}
	b.posToPiece[pos2] = piece
	delete(b.posToPiece, pos1)
	return nil
}

func (b *Board) Move(p1, p2 Pos) error {
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

	positions2 := getMovePositions(piece, p2)
	kingPos := b.kings[piece.Color^1]
	_, found = positions2[kingPos]
	if found && !b.moveBlocked(piece, p2, kingPos) {
		fmt.Println("CHECK")
	}

	// Move the piece to the new position.
	b.posToPiece[p2] = piece

	// Remove the piece from the old position.
	delete(b.posToPiece, p1)

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
	piece2, found := b.posToPiece[p2]

	// Check if piece2 is on the same team as piece.
	if found && piece.Color == piece2.Color {
		return ErrOccupiedPosition
	}

	// Check if the move from p1 to p2 is blocked by any other pieces.
	if b.moveBlocked(piece, p1, p2) {
		return ErrMoveBlocked
	}

	// See if move unblocks a path and now puts the king in check.
	//
	// Store below type moveBlocked positions at all times something
	// would be in check in map such as b.blockedChecks[Pos]*Piece.

	// If the piece is a King, see if by moving, it puts itself in check.
	if piece.Name == King {
		for pos, piece := range b.posToPiece {
			piecePositions := getMovePositions(piece, pos)
			_, checkFound := piecePositions[p2]
			// if checkFound && !b.moveBlocked(piece, pos, p2) {
			if checkFound {
				return ErrMovingIntoCheck
			}
		}
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
		if _, blocked := b.posToPiece[Pos{p1.X, p1.Y + d}]; blocked {
			return true
		}
	case Rook:
		switch {
		case p1.Y != p2.Y:
			for y := p1.Y + yd; y != p2.Y; y = y + yd {
				_, blocked := b.posToPiece[Pos{p2.X, y}]
				if blocked {
					return true
				}
			}
		case p1.X != p2.X:
			for x := p1.X + xd; x != p2.X; x = x + xd {
				_, blocked := b.posToPiece[Pos{x, p2.Y}]
				if blocked {
					return true
				}
			}
		}
	case Bishop:
		for x, y := p1.X+xd, p1.Y+yd; x != p2.X && y != p2.Y; x, y = x+xd, y+yd {
			// _, blocked := b.posToPiece[Pos{x, p2.Y}]
			// if blocked {
			// 	return true
			// }
		}
	}
	return false
}
