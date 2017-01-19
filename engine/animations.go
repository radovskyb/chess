package engine

import (
	"strings"
	"time"
)

// Animate takes a string of moves separated by commas and moves
// each piece in an animated fashion by sleeping for duration d
// in between printing the board after each move is made.
func Animate(b *Board, d time.Duration, movesStr string) error {
	moves := strings.Split(movesStr, ",")
	for _, move := range moves {
		if err := b.MoveByLocation(move[0:2], move[2:]); err != nil {
			return err
		}
		time.Sleep(d)
		b.Print()
	}
	return nil
}

// FourMoveCheckmate prints an animated four move checkmate game.
func FourMoveCheckmate(b *Board, d time.Duration) {
	Animate(b, d, "e2e4,e7e5,f1c4,b8c6,d1h5,d7d6,h5f7")
}
