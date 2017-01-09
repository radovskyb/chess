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

func diagPositions(cur Pos) map[Pos]struct{} {
	avail := make(map[Pos]struct{})
	for x, y := cur.X, cur.Y; x < 8 && y < 8; x, y = x+1, y+1 {
		avail[Pos{x, y}] = struct{}{}
	}
	for x, y := cur.X, cur.Y; x >= 0 && y < 8; x, y = x-1, y+1 {
		avail[Pos{x, y}] = struct{}{}
	}
	for x, y := cur.X, cur.Y; x < 8 && y >= 0; x, y = x+1, y-1 {
		avail[Pos{x, y}] = struct{}{}
	}
	for x, y := cur.X, cur.Y; x >= 0 && y >= 0; x, y = x-1, y-1 {
		avail[Pos{x, y}] = struct{}{}
	}
	return avail
}

func linePositions(cur Pos) map[Pos]struct{} {
	avail := make(map[Pos]struct{})
	for x := cur.X; x < 8; x++ {
		avail[Pos{x, cur.Y}] = struct{}{}
	}
	for x := cur.X; x >= 0; x-- {
		avail[Pos{x, cur.Y}] = struct{}{}
	}
	for y := cur.Y; y < 8; y++ {
		avail[Pos{cur.X, y}] = struct{}{}
	}
	for y := cur.Y; y >= 0; y-- {
		avail[Pos{cur.X, y}] = struct{}{}
	}
	return avail
}

func availablePositions(name PieceName, cur Pos) map[Pos]struct{} {
	avail := make(map[Pos]struct{})
	switch name {
	case Pawn:
		avail[Pos{cur.X, cur.Y + 1}] = struct{}{}
	case Knight:
		avail[Pos{cur.X + 2, cur.Y + 1}] = struct{}{}
		avail[Pos{cur.X - 2, cur.Y + 1}] = struct{}{}
		avail[Pos{cur.X + 2, cur.Y - 1}] = struct{}{}
		avail[Pos{cur.X - 2, cur.Y - 1}] = struct{}{}
		avail[Pos{cur.X + 1, cur.Y + 2}] = struct{}{}
		avail[Pos{cur.X - 1, cur.Y + 2}] = struct{}{}
		avail[Pos{cur.X + 1, cur.Y - 2}] = struct{}{}
		avail[Pos{cur.X - 1, cur.Y - 2}] = struct{}{}
	case Bishop:
		avail = diagPositions(cur)
	case Rook:
		avail = linePositions(cur)
	case Queen:
		avail = diagPositions(cur)
		for k, v := range linePositions(cur) {
			avail[k] = v
		}
	case King:
		avail[Pos{cur.X + 1, cur.Y}] = struct{}{}
		avail[Pos{cur.X - 1, cur.Y}] = struct{}{}
		avail[Pos{cur.X, cur.Y + 1}] = struct{}{}
		avail[Pos{cur.X, cur.Y - 1}] = struct{}{}
	}
	return avail
}
