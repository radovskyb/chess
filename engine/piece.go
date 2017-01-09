package engine

import "fmt"

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
	var color string
	switch p.Color {
	case Black:
		color = "b"
	case White:
		color = "w"
	default:
		return "invalid color"
	}
	name, found := pieceNames[p.Name]
	if !found {
		return "invalid piece name"
	}
	return fmt.Sprintf("%2s%s", color, name)
}
