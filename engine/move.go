package engine

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

// Move moves a piece on a board from positions p1 to p2.
//
// Move returns any errors that occur by trying to make
// the move from p1 to p2.
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

	// Get the move positions for the piece now at position p2.
	positions := getMovePositions(piece, p2)

	// Get the positions of the opponents king.
	kingPos := b.kings[piece.Color^1]

	// Check if the king's position is found within any of the
	// move positions for piece at p2.
	_, found = positions[kingPos]

	// If the king's position was found as isn't blocked, it's a check.
	if found && !b.moveBlocked(piece, p2, kingPos) {
		b.check[piece.Color^1] = true
	}

	// Move the piece to the new position.
	b.posToPiece[p2] = piece

	// Remove the piece from the old position.
	delete(b.posToPiece, p1)

	// Update current king's position.
	if piece.Name == King {
		b.kings[piece.Color] = p2
	}

	// Update who's turn it is.
	b.Turn ^= 1

	return nil
}

// moveLegal checks to see whether the specified move is legal to make or not.
func (b *Board) moveLegal(piece *Piece, p1, p2 Pos) error {
	// Get a list of all possible move positions that the
	// piece can move to without restrictions.
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

	// Pawn is moving yd+sideways, make sure there's an opponents piece at p2.
	if !found && piece.Name == Pawn && p1.X != p2.X {
		return ErrInvalidPieceMove
	}

	// Check if the move from p1 to p2 is blocked by any other pieces.
	if b.moveBlocked(piece, p1, p2) {
		return ErrMoveBlocked
	}

	// If the piece is a King, see if by moving, it puts itself in check.
	if piece.Name == King {
		// Remove the king from it's current position on the board
		// so that it doesn't block any pieces in position attacked.
		delete(b.posToPiece, p1)

		// Check if the position is being attacked
		attacked := b.positionAttacked(p2, piece.Color)

		// Put the king back at position p1.
		b.posToPiece[p1] = piece

		// If the position is being attacked, the move is illegal
		// as the king would be moving into a check.
		if attacked {
			return ErrMovingIntoCheck
		}
	}

	return nil
}

func (b *Board) positionAttacked(at Pos, c Color) bool {
	for pos, piece := range b.posToPiece {
		if piece.Color == c {
			continue
		}
		positions := getMovePositions(piece, pos)
		_, found := positions[at]
		if found && !b.moveBlocked(piece, pos, at) {
			return true
		}
	}
	return false
}

func (b *Board) diagMoveBlocked(p1, p2 Pos, xd, yd int) bool {
	for x, y := p1.X+xd, p1.Y+yd; x != p2.X && y != p2.Y; x, y = x+xd, y+yd {
		_, blocked := b.posToPiece[Pos{x, y}]
		if blocked {
			return true
		}
	}
	return false
}

func (b *Board) lineMoveBlocked(p1, p2 Pos, xd, yd int) bool {
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
	return false
}

func (b *Board) moveBlocked(piece *Piece, p1, p2 Pos) bool {
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
		if p1.X != p2.X {
			return false
		}
		if _, blocked := b.posToPiece[Pos{p1.X, p1.Y + d}]; blocked {
			return true
		}
	case Rook:
		return b.lineMoveBlocked(p1, p2, xd, yd)
	case Bishop:
		return b.diagMoveBlocked(p1, p2, xd, yd)
	case Queen:
		if p1.Y == p2.Y || p1.X == p2.X {
			return b.lineMoveBlocked(p1, p2, xd, yd)
		}
		return b.diagMoveBlocked(p1, p2, xd, yd)
	}
	return false
}
