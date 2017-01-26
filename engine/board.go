package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	ErrInvalidLocation        = errors.New("error: location string is invalid")
	ErrNoPieceAtPosition      = errors.New("error: no piece at specified position")
	ErrOpponentsPiece         = errors.New("error: piece belongs to opponent")
	ErrInvalidPieceMove       = errors.New("error: invalid move for piece")
	ErrOccupiedPosition       = errors.New("error: position is already occupied")
	ErrMoveBlocked            = errors.New("error: move is blocked by another piece")
	ErrMovingIntoCheck        = errors.New("error: move puts king in check")
	ErrMoveWhileInCheck       = errors.New("error: can't move piece while in check")
	ErrNoRookToCastleWith     = errors.New("error: rook not found to castle with")
	ErrKingOrRookMoved        = errors.New("error: king or rook has already moved before")
	ErrCastleWithKingInCheck  = errors.New("error: castle while king is in check")
	ErrCastleWithPieceBetween = errors.New("error: castle with pieces between king and rook")
	ErrCastleMoveThroughCheck = errors.New("error: castle moving king through check")
	ErrNoPreviousMove         = errors.New("error: no previous move available")
	ErrKingTooCloseToKing     = errors.New("error: king can't be that close to another king")
)

type Color uint8

const (
	Black Color = iota
	White
)

// String returns a string for a Color.
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

// piecePos contains a *Piece and it's position on the board.
type piecePos struct {
	*Piece
	Pos
}

// A Board describes a chess board.
type Board struct {
	// turn holds a color value for who's turn it is.
	turn Color

	// posToPiece holds a map of positions to pieces on the board.
	posToPiece map[Pos]*Piece

	// check holds a true or false based on whether either king
	// is currently in check or not.
	check [2]bool

	// kings holds both of the king's positions on the board.
	kings [2]Pos

	// kingLos holds pieces that have a line of sight to a king.
	kingLos [2]map[piecePos]struct{}

	// history contains a slice of all moves that have occurred
	// on the board.
	history []*move

	// moveNum stores the current move index in the history slice.
	moveNum int

	// hasMoved holds any pieces that have already moved in the game.
	hasMoved map[*Piece]int

	// mustPromote holds which color needs to promote a pawn.
	mustPromote [2]bool
}

func (b *Board) Turn() Color {
	return b.turn
}

// History returns a string combining all of the moves currently
// stored in the board's history in the format of l1l2,l1l2 etc.
func (b *Board) History() (history string) {
	for _, m := range b.history {
		history += fmt.Sprintf("%s%s,", m.from, m.to)
	}
	return strings.TrimRight(history, ",")
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

// MustPromote returns true and the color of the side that must
// promote a pawn or false if no piece currently needs to promote.
func (b *Board) MustPromote() (bool, Color) {
	if b.mustPromote[Black] {
		return true, Black
	}
	if b.mustPromote[White] {
		return true, White
	}
	return false, 2 // 2 for invalid color.
}

// NewBoard creates an initializes a new chess board.
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
	// Create and initalize a new map for b.hasMoved.
	hasMoved := make(map[*Piece]int)
	for _, piece := range posToPiece {
		hasMoved[piece] = 0
	}
	return &Board{
		turn:       White,
		posToPiece: posToPiece,
		kings:      [2]Pos{White: {4, 0}, Black: {4, 7}},
		kingLos:    [2]map[piecePos]struct{}{White: {}, Black: {}},
		history:    []*move{}, // Create a new blank history.
		moveNum:    -1,
		hasMoved:   hasMoved,
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

// Print prints the board in terminals using ansi escape codes and colors.
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
