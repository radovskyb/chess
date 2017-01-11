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

var pieceNames = map[PieceName]string{
	Pawn: "P", Knight: "N", Bishop: "B",
	Rook: "R", Queen: "Q", King: "K",
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
