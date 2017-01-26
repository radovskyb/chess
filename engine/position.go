package engine

type Pos struct {
	X, Y int
}

// String returns a Pos as a location string. Example: Pos{0, 0} -> A1.
func (p Pos) String() string {
	return string(p.X+'A') + string(p.Y+'1')
}

// locToPos turns a location string into a Pos object.
// If the location is invalid, an error and an invalid position is returned.
func locToPos(loc string) (Pos, error) {
	if len(loc) != 2 {
		return Pos{-1, -1}, ErrInvalidLocation
	}
	if !(loc[0] < 'a' || loc[0] > 'h' || loc[1] < '1' || loc[1] > '8') {
		return Pos{int(loc[0] - 'a'), int(loc[1] - '1')}, nil
	}
	if !(loc[0] < 'A' || loc[0] > 'H' || loc[1] < '1' || loc[1] > '8') {
		return Pos{int(loc[0] - 'A'), int(loc[1] - '1')}, nil
	}
	return Pos{-1, -1}, ErrInvalidLocation
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
			if cur.Y != 0 {
				pos[Pos{cur.X, cur.Y - 1}] = struct{}{}
				if cur.X != 0 {
					pos[Pos{cur.X - 1, cur.Y - 1}] = struct{}{}
				}
				if cur.X != 7 {
					pos[Pos{cur.X + 1, cur.Y - 1}] = struct{}{}
				}
			}
		case White:
			if cur.Y == 1 {
				pos[Pos{cur.X, cur.Y + 2}] = struct{}{}
			}
			if cur.Y != 7 {
				pos[Pos{cur.X, cur.Y + 1}] = struct{}{}
				if cur.X != 0 {
					pos[Pos{cur.X - 1, cur.Y + 1}] = struct{}{}
				}
				if cur.X != 7 {
					pos[Pos{cur.X + 1, cur.Y + 1}] = struct{}{}
				}
			}
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

func (b *Board) availBetween(piece *Piece, p1, p2 Pos) map[Pos]struct{} {
	yd, xd := 1, 1
	if p1.Y > p2.Y {
		yd = -1
	}
	if p1.X > p2.X {
		xd = -1
	}
	switch piece.Name {
	case Pawn:
		d := 1
		if piece.Color == Black {
			d = -1
		}
		if p1.Y == p2.Y+(2*d) {
			if _, blocked := b.posToPiece[Pos{p1.X, p1.Y + d}]; blocked {
				return map[Pos]struct{}{Pos{p1.X, p1.Y + d}: struct{}{}}
			}
		}
	case Rook:
		return b.lineBetween(p1, p2, xd, yd)
	case Bishop:
		return b.diagBetween(p1, p2, xd, yd)
	case Queen:
		if p1.Y == p2.Y || p1.X == p2.X {
			return b.lineBetween(p1, p2, xd, yd)
		}
		return b.diagBetween(p1, p2, xd, yd)
	}
	return nil
}

func (b *Board) lineBetween(p1, p2 Pos, xd, yd int) map[Pos]struct{} {
	m := map[Pos]struct{}{}
	switch {
	case p1.Y != p2.Y:
		for y := p1.Y + yd; y != p2.Y; y = y + yd {
			_, blocked := b.posToPiece[Pos{p2.X, y}]
			if !blocked {
				m[Pos{p2.X, y}] = struct{}{}
			}
		}
	case p1.X != p2.X:
		for x := p1.X + xd; x != p2.X; x = x + xd {
			_, blocked := b.posToPiece[Pos{x, p2.Y}]
			if !blocked {
				m[Pos{x, p2.Y}] = struct{}{}
			}
		}
	}
	return m
}

func (b *Board) diagBetween(p1, p2 Pos, xd, yd int) map[Pos]struct{} {
	m := map[Pos]struct{}{}
	for x, y := p1.X+xd, p1.Y+yd; x != p2.X && y != p2.Y; x, y = x+xd, y+yd {
		_, blocked := b.posToPiece[Pos{x, y}]
		if !blocked {
			m[Pos{x, y}] = struct{}{}
		}
	}
	return m
}

func (b *Board) positionOffBoard(pos Pos) bool {
	if pos.X < 0 || pos.X > 7 || pos.Y < 0 || pos.Y > 7 {
		return true
	}
	return false
}
