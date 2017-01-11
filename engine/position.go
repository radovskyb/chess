package engine

import "strings"

type Pos struct {
	X, Y int
}

var letterToX = map[string]int{
	"a": 0, "b": 1, "c": 2, "d": 3,
	"e": 4, "f": 5, "g": 6, "h": 7,
}
var letterToY = map[string]int{
	"1": 0, "2": 1, "3": 2, "4": 3,
	"5": 4, "6": 5, "7": 6, "8": 7,
}

// locToPos turns a location string into a Pos object.
// If the location is invalid, an error is returned and an empty Pos object is returned.
//
// Since the empty Pos is still a valid position (x: 0, y: 0), the error must be checked
// before using Pos after a call to locToPos.
func locToPos(loc string) (Pos, error) {
	loc = strings.ToLower(loc)
	if len(loc) != 2 {
		return Pos{}, ErrInvalidLocation
	}
	x, foundX := letterToX[string(loc[0])]
	y, foundY := letterToY[string(loc[1])]
	if !(foundX && foundY) {
		return Pos{}, ErrInvalidLocation
	}
	return Pos{x, y}, nil
}

// diagPositions returns a map of diagonal move positions starting
// from the specified current position.
func diagPositions(cur Pos) map[Pos]struct{} {
	pos := make(map[Pos]struct{})
	for x, y := cur.X+1, cur.Y+1; x < 8 && y < 8; x, y = x+1, y+1 {
		pos[Pos{x, y}] = struct{}{}
	}
	for x, y := cur.X-1, cur.Y+1; x >= 0 && y < 8; x, y = x-1, y+1 {
		pos[Pos{x, y}] = struct{}{}
	}
	for x, y := cur.X+1, cur.Y-1; x < 8 && y >= 0; x, y = x+1, y-1 {
		pos[Pos{x, y}] = struct{}{}
	}
	for x, y := cur.X-1, cur.Y-1; x >= 0 && y >= 0; x, y = x-1, y-1 {
		pos[Pos{x, y}] = struct{}{}
	}
	return pos
}

// linePositions returns a map of straight line move positions starting
// from the specified current position.
func linePositions(cur Pos) map[Pos]struct{} {
	pos := make(map[Pos]struct{})
	for x := cur.X + 1; x < 8; x++ {
		pos[Pos{x, cur.Y}] = struct{}{}
	}
	for x := cur.X - 1; x >= 0; x-- {
		pos[Pos{x, cur.Y}] = struct{}{}
	}
	for y := cur.Y + 1; y < 8; y++ {
		pos[Pos{cur.X, y}] = struct{}{}
	}
	for y := cur.Y - 1; y >= 0; y-- {
		pos[Pos{cur.X, y}] = struct{}{}
	}
	return pos
}

// getMovePositions returns a map of all possible positions that the
// specified piece could move to with no restrictions in place.
func getMovePositions(piece *Piece, cur Pos) map[Pos]struct{} {
	pos := make(map[Pos]struct{})

	switch piece.Name {
	case Pawn:
		switch piece.Color {
		case Black:
			if cur.Y == 6 {
				pos[Pos{cur.X, cur.Y - 2}] = struct{}{}
			}
			pos[Pos{cur.X, cur.Y - 1}] = struct{}{}
		case White:
			if cur.Y == 1 {
				pos[Pos{cur.X, cur.Y + 2}] = struct{}{}
			}
			pos[Pos{cur.X, cur.Y + 1}] = struct{}{}
		}
	case Knight:
		pos[Pos{cur.X + 2, cur.Y + 1}] = struct{}{}
		pos[Pos{cur.X - 2, cur.Y + 1}] = struct{}{}
		pos[Pos{cur.X + 2, cur.Y - 1}] = struct{}{}
		pos[Pos{cur.X - 2, cur.Y - 1}] = struct{}{}
		pos[Pos{cur.X + 1, cur.Y + 2}] = struct{}{}
		pos[Pos{cur.X - 1, cur.Y + 2}] = struct{}{}
		pos[Pos{cur.X + 1, cur.Y - 2}] = struct{}{}
		pos[Pos{cur.X - 1, cur.Y - 2}] = struct{}{}
	case Bishop:
		pos = diagPositions(cur)
	case Rook:
		pos = linePositions(cur)
	case Queen:
		pos = diagPositions(cur)
		for k, v := range linePositions(cur) {
			pos[k] = v
		}
	case King:
		pos[Pos{cur.X + 1, cur.Y}] = struct{}{}
		pos[Pos{cur.X - 1, cur.Y}] = struct{}{}
		pos[Pos{cur.X, cur.Y + 1}] = struct{}{}
		pos[Pos{cur.X, cur.Y - 1}] = struct{}{}
		pos[Pos{cur.X + 1, cur.Y + 1}] = struct{}{}
		pos[Pos{cur.X + 1, cur.Y - 1}] = struct{}{}
		pos[Pos{cur.X - 1, cur.Y + 1}] = struct{}{}
		pos[Pos{cur.X - 1, cur.Y - 1}] = struct{}{}
	}
	return pos
}
