package engine

import "fmt"

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
	if piece.Color != b.turn {
		return ErrOpponentsPiece
	}

	// Castling.
	if piece.Name == King && p1.X == 4 &&
		(p2.Y == 0 || p2.Y == 7) &&
		(p2.X == 2 || p2.X == 6) {
		return b.doCastling(piece, p1, p2)
	}

	// Check if the move is legal to make.
	if err := b.moveLegal(piece, p1, p2); err != nil {
		return err
	}

	// Remove the piece from the old position p1 here, so it
	// doesn't block when checking b.moveBlocked below if waiting
	// to delete when adding piece to p2.
	delete(b.posToPiece, p1)

	// Get the move positions for the piece now at position p2.
	positions := getMovePositions(piece, p2)

	// Get the positions of the opponents king.
	kingPos := b.kings[piece.Color^1]

	// Check if the king's position is found within any of the
	// move positions for piece at p2.
	_, found = positions[kingPos]

	// If the king's position was found as isn't blocked, it's a check.
	if found {
		if !b.moveBlocked(piece, p2, kingPos) {
			b.check[piece.Color^1] = true
		}
		b.kingLos[piece.Color^1] = append(b.kingLos[piece.Color^1],
			piecePos{piece, p2})
	}

	// Move the piece to the new position.
	b.posToPiece[p2] = piece

	// Update current king's position.
	if piece.Name == King {
		b.kings[piece.Color] = p2
	}

	// If color's king was in check and the current move
	// is legal, the king will no longer be in check.
	b.check[piece.Color] = false

	// Add piece to the hasMoved map.
	b.hasMoved[piece] = struct{}{}

	// Update who's turn it is.
	b.turn ^= 1

	return nil
}

// moveLegal checks to see whether the specified move is legal to
// make or not.
func (b *Board) moveLegal(piece *Piece, p1, p2 Pos) error {
	// Get a list of all possible move positions that the
	// piece can move to without restrictions.
	positions := getMovePositions(piece, p1)
	if _, ok := positions[p2]; !ok {
		return ErrInvalidPieceMove
	}

	// Check if there's a piece at position p2.
	piece2, found := b.posToPiece[p2]

	// If there was a piece found at position p2.
	if found {
		// If piece2 is the color as piece, the position
		// is occupied.
		if piece.Color == piece2.Color {
			return ErrOccupiedPosition
		}
		// If the piece is a pawn, it can't take a piece in
		// front of it.
		if piece.Name == Pawn && p1.X == p2.X {
			return ErrOccupiedPosition
		}
	}

	// Check if the move from p1 to p2 is blocked by any other pieces.
	if b.moveBlocked(piece, p1, p2) {
		return ErrMoveBlocked
	}

	// If color is in check, make sure that the piece can't
	// be moved, unless it evades the current check.
	//
	// If not in check, then make sure move doesn't put own king in check.
	if piece.Name != King {
		for i, pp := range b.kingLos[piece.Color] {
			// Check if piece is still at pp.Pos. If it isn't, delete
			// pp.Piece from the kings line of sight slice for color.
			pc, found := b.posToPiece[pp.Pos]
			if !found || pc != pp.Piece {
				// Delete from b.kingLos[piece.Color]
				b.kingLos[piece.Color] = append(
					b.kingLos[piece.Color][:i],
					b.kingLos[piece.Color][i:]...,
				)
				continue
			}

			// If piece is trying to take pp.Piece, continue.
			if p2 == pp.Pos {
				continue
			}

			// If pp.Piece is not still attacking the king, delete it
			// from the b.kingLos[piece.Color] slice.
			positions := getMovePositions(pp.Piece, pp.Pos)
			if _, kingFound := positions[b.kings[piece.Color]]; !kingFound {
				// Delete from b.kingLos[piece.Color]
				b.kingLos[piece.Color] = append(
					b.kingLos[piece.Color][:i],
					b.kingLos[piece.Color][i:]...,
				)
				continue
			}

			// If the piece isn't in pp.Piece's line of sight and
			// it's not currently a check, the piece can move since
			// it won't open up any new checks.
			_, p1found := positions[p1]
			if !b.check[piece.Color] && !p1found {
				continue
			}

			// Simulate moving piece to p2 and then check if pp.Piece is
			// blocked to the king or not.
			//
			// Get the piece at p2.
			pc2 := b.posToPiece[p2]
			// Move p1's piece to p2.
			b.posToPiece[p2] = piece
			// Delete piece from p1.
			delete(b.posToPiece, p1)
			// Check if pp.Piece is blocked to the king.
			blocked := b.moveBlocked(pp.Piece, pp.Pos, b.kings[piece.Color])
			// Move p1's piece back to p1.
			b.posToPiece[p1] = piece
			// Delete p1's piece from p2.
			delete(b.posToPiece, p2)
			// If p2 originally contained a piece, put it back.
			if pc2 != nil {
				b.posToPiece[p2] = pc2
			}

			if !blocked {
				if b.check[piece.Color] {
					// If there would still be a check, return an
					// ErrMoveWhileInCheck error.
					return ErrMoveWhileInCheck
				}
				return ErrMovingIntoCheck
			}
		}
	}

	// If the piece is a King, see if it can move to p2 without
	// being put into check.
	if piece.Name == King {
		// Remove the king from it's current position on the board
		// so that it doesn't block any pieces in position attacked.
		delete(b.posToPiece, p1)

		// Check if the position is being attacked by the opponent's color.
		attacked := b.positionAttacked(p2, piece.Color^1)

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

// positionAttacked returns a true or false based on whether the
// at position is being attacking by any pieces from color by.
func (b *Board) positionAttacked(at Pos, by Color) bool {
	for pos, piece := range b.posToPiece {
		if piece.Color != by || (piece.Name == Pawn && pos.X == at.X) {
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

// moveBlocked checks whether there are any pieces inbetween piece
// at position p1 and anything at position p2.
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
		// TODO: Fix for En Passant.
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

func (b *Board) doCastling(king *Piece, p1, p2 Pos) error {
	if b.check[king.Color] {
		return ErrCastleWithKingInCheck
	}

	if _, found := b.hasMoved[king]; found {
		return ErrKingOrRookMoved
	}

	switch p2.X {
	case 2: // Queen-side.
		// Move the queen-side rook to d1 or d8.
		piece, found := b.posToPiece[Pos{0, p2.Y}]
		if !found {
			return ErrNoRookToCastleWith
		}
		if _, found := b.hasMoved[piece]; found {
			return ErrKingOrRookMoved
		}

		// Make sure there's no pieces in between the king and the rook.
		for x := 1; x < 4; x++ {
			if _, found := b.posToPiece[Pos{x, p2.Y}]; found {
				return ErrCastleWithPieceBetween
			}
		}

		for x := 2; x < 4; x++ {
			if b.positionAttacked(Pos{x, p2.Y}, piece.Color^1) {
				return ErrCastleMoveThroughCheck
			}
		}

		// Add the rook to it's new postion.
		b.posToPiece[Pos{3, p2.Y}] = piece

		// Remove the rook from the old position.
		delete(b.posToPiece, Pos{0, p2.Y})
	case 6: // King-side.
		// Move the king-side rook to f1 or f8.
		piece, found := b.posToPiece[Pos{7, p2.Y}]
		if !found {
			return ErrNoRookToCastleWith
		}
		if _, found := b.hasMoved[piece]; found {
			return ErrKingOrRookMoved
		}

		// Make sure there's no pieces in between the king and the rook.
		for x := 5; x < 7; x++ {
			if _, found := b.posToPiece[Pos{x, p2.Y}]; found {
				return ErrCastleWithPieceBetween
			}
			if b.positionAttacked(Pos{x, p2.Y}, piece.Color^1) {
				return ErrCastleMoveThroughCheck
			}
		}

		// Add the rook to it's new postion.
		b.posToPiece[Pos{5, p2.Y}] = piece

		// Remove the rook from the old position.
		delete(b.posToPiece, Pos{7, p2.Y})
	default:
		return fmt.Errorf("cant castle king to position %s", p2)
	}

	// Move the king to p2.
	b.posToPiece[p2] = king
	// Remove the piece from the old position.
	delete(b.posToPiece, p1)
	// Update current king's position.
	b.kings[king.Color] = p2
	// Add king to the hasMoved map. The rook doesn't need to be
	// added to the map since a king can only castle once.
	b.hasMoved[king] = struct{}{}
	// Update who's turn it is.
	b.turn ^= 1

	return nil
}
