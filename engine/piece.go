package engine

import "github.com/fatih/color"

type PieceName uint8

const (
	Pawn PieceName = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

func (n PieceName) String() string {
	switch n {
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	default:
		return "invalid piece name"
	}
}

var pieceNames = map[PieceName]string{
	Pawn: "\u2659", Knight: "\u2658", Bishop: "\u2657",
	Rook: "\u2656", Queen: "\u2655", King: "\u2654",
}

type Piece struct {
	Name PieceName
	Color
}

func (p *Piece) String() string {
	name, found := pieceNames[p.Name]
	if !found {
		return "invalid piece name"
	}
	switch p.Color {
	case Black:
		return color.BlackString(" %s ", name)
	case White:
		return color.WhiteString(" %s ", name)
	}
	return "invalid color"
}
