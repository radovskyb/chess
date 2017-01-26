package engine

import "fmt"

type move struct {
	piece     *Piece
	from      Pos
	to        Pos
	captured  *Piece
	enPassant bool // whether or not the capture was en passant.
}

func (b *Board) makeMove(m *move) {
	// Remove the piece from the old position from over here, so it
	// doesn't block when checking b.moveBlocked below if waiting
	// to delete when adding piece to position to.
	delete(b.posToPiece, m.from)

	// Move the piece to the new position.
	b.posToPiece[m.to] = m.piece

	// If the move is an en passant, delete the captured
	// piece from the board.
	if m.enPassant {
		delete(b.posToPiece, Pos{m.to.X, m.from.Y})
	}

	// Update current king's position and line of sights.
	if m.piece.Name == King {
		b.kings[m.piece.Color] = m.to

		// If the king is moving into any opponent's piece's line
		// of sight, add it to the b.kingLos slice.
		for pos, piece := range b.posToPiece {
			if m.piece.Color != piece.Color^1 || (piece.Name == Pawn && pos.X == m.to.X) {
				continue
			}
			positions := getMovePositions(piece, pos)
			_, found := positions[m.to]
			// Don't need to check if b.moveBlocked since already checked
			// in b.moveLegal.
			if found {
				b.kingLos[m.piece.Color][piecePos{piece, pos}] = struct{}{}
			}
		}
	}

	// Get the move positions for the piece now at position to.
	positions := getMovePositions(m.piece, m.to)

	// Get the positions of the opponent's king.
	kingPos := b.kings[m.piece.Color^1]

	// Check if the king's position is found within any of the
	// move positions for piece at position to.
	_, found := positions[kingPos]

	// If the king's position was found as isn't blocked, it's a check.
	if found {
		// If the piece at it's new position is now causing the opponent's
		// king to be in check, set b.check to true for the color.
		if !b.moveBlocked(m.piece, m.to, kingPos) {
			b.check[m.piece.Color^1] = true
		}
		b.kingLos[m.piece.Color^1][piecePos{m.piece, m.to}] = struct{}{}
	}

	// Increment b.hasMoved for piece.
	b.hasMoved[m.piece]++

	// If the history's length has already reached b.moveNum, it means
	// that the previous move was an undo and since this new move will now
	// create a new forward history, delete any moves after this move in the
	// history slice, since they're no longer relevant history.
	if len(b.history) > b.moveNum {
		b.history = b.history[:b.moveNum+1]
	}

	// Increment the history move index number.
	b.moveNum++

	// Add the new move to b.history.
	b.history = append(b.history, m)

	// Update who's turn it is.
	b.turn ^= 1
}

// MoveByLocation is a convenience method that makes a move based
// 2 location strings instead of Pos objects. For example, a2 to a4.
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

// newMove creates a new move.
func (b *Board) newMove(piece *Piece, from, to Pos, enPassant bool) *move {
	m := &move{
		piece:     piece,
		from:      from,
		to:        to,
		enPassant: enPassant,
	}
	if enPassant {
		m.captured = b.posToPiece[Pos{to.X, from.Y}]
	} else {
		m.captured = b.posToPiece[to]
	}
	return m
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

	// Make sure that p2 is a valid move position for piece.
	positions := getMovePositions(piece, p1)
	if _, ok := positions[p2]; !ok {
		return ErrInvalidPieceMove
	}

	if err := b.moveLegal(piece, p1, p2); err != nil {
		return err
	}

	// Make the move on the board.
	//
	// If the move is legal to make and the piece is a pawn
	// trying to move diagonally, but there's no piece at
	// position p2, it's an en passant, otherwise make a
	// normal move.
	if piece.Name == Pawn && p1.X != p2.X {
		_, found := b.posToPiece[p2]
		if found {
			b.makeMove(b.newMove(piece, p1, p2, false))
		} else {
			b.makeMove(b.newMove(piece, p1, p2, true))
		}
	} else {
		b.makeMove(b.newMove(piece, p1, p2, false))
	}

	// If color's king was in check and the current move
	// is legal, the king will no longer be in check.
	b.check[piece.Color] = false

	return nil
}

// InCheckmate returns a true or false based on whether the
// color is currently in checkmate or not.
func (b *Board) InCheckmate(color Color) bool {
	// If the king can move, it's not a checkmate.
	if b.kingCanMove(color) {
		return false
	}
	// If a piece can be made to stop all checks.
	if b.canStopAllChecks(color) {
		return false
	}
	return true
}

// HasStalemate checks if there's currently a stalemate on the
// board.
func (b *Board) HasStalemate(color Color) bool {
	// If the king can move, it's not a stalemate.
	if b.kingCanMove(color) {
		return false
	}

	// If any piece can move, it's not a stalemate.
	if b.anyPieceCanMove(color) {
		return false
	}

	// There is a stalemate on the board as color currently
	// has nowhere legal to move on the board.
	return true
}

func (b *Board) anyPieceCanMove(color Color) bool {
	for pos, pc := range b.posToPiece {
		if pc.Color != color || pc.Name == King {
			continue
		}
		// Get all possible positions for piece at pos.
		positions := getMovePositions(pc, pos)
		// Iterate over every position that pc can move to and if
		// pc can move to any position, return true.
		for p := range positions {
			if b.moveLegal(pc, pos, p) == nil {
				return true
			}
		}
	}
	return false
}

func (b *Board) canStopAllChecks(color Color) bool {
	// Get all blocking positions between or on the line
	// of sights and the king for colorq.
	betweenOrOn := make(map[Pos]struct{})

	// Iterate over all pieces in king's line of sight.
	for pp := range b.kingLos[color] {
		// Add the piece's position to the betweenOrOn map.
		betweenOrOn[pp.Pos] = struct{}{}

		// Add all of the pieces between the king and the piece
		// to the betweenOrOn map.
		for k, v := range b.availBetween(pp.Piece, pp.Pos, b.kings[color]) {
			betweenOrOn[k] = v
		}
	}

	// Try to move to all positions that any white piece shares with
	// positions in the between map and if any positions prevent the king
	// from being in check, it's no longer a checkmate.
	//
	// Iterate over all pieces for color.
	for pos, pc := range b.posToPiece {
		if pc.Color != color || pc.Name == King {
			continue
		}

		// Get all possible positions for piece at pos.
		positions := getMovePositions(pc, pos)

		// Iterate over every position that pc can move to.
		for p := range positions {
			// If position p doesn't exist in map between, continue.
			if _, ok := betweenOrOn[p]; !ok {
				continue
			}

			// If piece can't move from pos to p, continue.
			if b.moveLegal(pc, pos, p) != nil {
				continue
			}

			// Simulate moving pc to p and then check if the
			// king is still in check.
			//
			// Get the piece at p.
			pc2 := b.posToPiece[p]
			// Move pc to p.
			b.posToPiece[p] = pc
			// Delete pc from pos.
			delete(b.posToPiece, pos)

			// If the king's no longer in check, put
			// the pieces back and return false.
			if !b.kingInCheck(color) {
				// Move piece back to pos.
				b.posToPiece[pos] = pc
				// Delete pc from p.
				delete(b.posToPiece, p)
				// If p originally contained a piece, put it back.
				if pc2 != nil {
					b.posToPiece[p] = pc2
				}
				return true
			}

			// Move piece back to pos.
			b.posToPiece[pos] = pc
			// Delete pc from p.
			delete(b.posToPiece, p)
			// If p originally contained a piece, put it back.
			if pc2 != nil {
				b.posToPiece[p] = pc2
			}

		}
	}

	return false
}

// kingInCheck is a wrapper around b.positionAttacked
// to see if the king for color's position is currently
// being attacked which would mean the king is in check.
func (b *Board) kingInCheck(color Color) bool {
	return b.positionAttacked(b.kings[color], color^1)
}

// kingCanMove determines whether the king for color
// has any positions that it can legally move to or not.
func (b *Board) kingCanMove(color Color) bool {
	// Get the positions of the opponent's king.
	kingPos := b.kings[color]
	// Get the king's piece.
	king := b.posToPiece[kingPos]
	// Get the king's potential move positions.
	kingPositions := getMovePositions(king, kingPos)
	// If there's any positions that the king can
	// legally move to, return true since the king
	// can move.
	for pos := range kingPositions {
		if b.moveLegal(king, kingPos, pos) == nil {
			return true
		}
	}
	// There was no legal moves for the king, return false.
	return false
}

// moveLegal checks to see whether the specified move is legal to
// make or not.
//
// moveLegal doesn't check if p2 is a possible available move for
// the piece type and should be checked before calling moveLegal.
func (b *Board) moveLegal(piece *Piece, p1, p2 Pos) error {
	if b.positionOffBoard(p2) {
		return fmt.Errorf("position p2 is off the board")
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

	// If there was no piece at position p2 and the piece is a
	// pawn trying to move diagonally, if there's no setup for
	// an en passant, the move is illegal.
	if !found && piece.Name == Pawn && p1.X != p2.X &&
		!b.canEnPassant(piece, p1, p2) {
		return ErrInvalidPieceMove
	}

	// Check if the move from p1 to p2 is blocked by any other pieces.
	if b.moveBlocked(piece, p1, p2) {
		return ErrMoveBlocked
	}

	// If color is in check, make sure that the piece can't
	// be moved, unless it evades the current check.
	//
	// If not in check, then make sure move doesn't put own king in check.
	if err := b.moveIntoOrWhileCheck(piece, p1, p2); err != nil {
		return err
	}

	// If the piece is a King, see if it can move to p2 without
	// being put into check and also make sure that the opponent's
	// king is not within 1 X or Y position of position p2.
	if piece.Name == King {
		// Get the king's possible move positions from position p2.
		kingPositions := getMovePositions(piece, p2)
		// Iterate over all the possible positions from p2 and see if
		// the opponent's king is in any of those positions.
		for pos := range kingPositions {
			if pc, found := b.posToPiece[pos]; found {
				// The king will be too close to the other king if it
				// moves to position p2.
				if pc.Name == King && pc.Color != piece.Color {
					return ErrKingTooCloseToKing
				}
			}
		}

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

// pieceCausingCheckTo returns a true or false based on whether the piece
// from piecePos pp causes the specified color's king to be in check.
func (b *Board) pieceCausingCheckTo(pp piecePos, color Color) bool {
	// Check if piece is still at pp.Pos. If it isn't, delete
	// pp.Piece from the kings line of sight slice for color.
	pc, found := b.posToPiece[pp.Pos]
	if !found || pc != pp.Piece {
		// Delete from b.kingLos[color]
		delete(b.kingLos[color], pp)
		return false
	}

	// If pp.Piece is not still attacking the king, delete it
	// from the b.kingLos[piece.Color] slice.
	positions := getMovePositions(pp.Piece, pp.Pos)
	if _, kingFound := positions[b.kings[color]]; !kingFound {
		// Delete from b.kingLos[piece.Color]
		delete(b.kingLos[color], pp)
		return false
	}

	if !b.moveBlocked(pp.Piece, pp.Pos, b.kings[color]) {
		return true
	}

	return false
}

// moveIntoOrWhileCheck returns an error if moving piece causes the piece's color's
// king to be in check, or if the king will still be in check if the move does not
// uncheck the king through capture or blockage if the king is already in check.
func (b *Board) moveIntoOrWhileCheck(piece *Piece, p1, p2 Pos) error {
	if piece.Name != King {
		for pp := range b.kingLos[piece.Color] {
			// Check if piece is still at pp.Pos. If it isn't, delete
			// pp.Piece from the kings line of sight slice for color.
			pc, found := b.posToPiece[pp.Pos]
			if !found || pc != pp.Piece {
				// Delete from b.kingLos[piece.Color]
				delete(b.kingLos[piece.Color], pp)
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
				delete(b.kingLos[piece.Color], pp)
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
	return nil
}

func (b *Board) canEnPassant(piece *Piece, p1, p2 Pos) bool {
	switch piece.Color {
	case Black:
		if p1.Y != 3 {
			return false
		}
	case White:
		if p1.Y != 4 {
			return false
		}
	}
	pc, ok := b.posToPiece[Pos{p2.X, p1.Y}]
	if ok && pc.Name == Pawn && pc.Color != piece.Color {
		// If the previous move on the board is not piece pc
		// moving from position Pos{p2.X, p1.Y - 2} for white
		// or Pos{p2.X, p1.Y + 2} for black to it's current
		// position at Pos{p2.X, p1.Y}, en passant is not allowed.
		prevMove, err := b.prevMove()
		if err != nil {
			return false
		}
		d := 1
		if piece.Color == Black {
			d = -1
		}
		if prevMove.piece != pc || prevMove.from.X != p2.X ||
			prevMove.to.X != p2.X || prevMove.from.Y != p1.Y+2*d {
			return false
		}

		// Pawn can en passant.
		return true
	}
	return false
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

// diagMoveBlocked checks whether there are any pieces between
// positions p1 and p2 in a diagonal line direction.
func (b *Board) diagMoveBlocked(p1, p2 Pos, xd, yd int) bool {
	for x, y := p1.X+xd, p1.Y+yd; x != p2.X && y != p2.Y; x, y = x+xd, y+yd {
		_, blocked := b.posToPiece[Pos{x, y}]
		if blocked {
			return true
		}
	}
	return false
}

// lineMoveBlocked checks whether there are any pieces between
// positions p1 and p2 in a straight line direction.
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

// moveBlocked checks whether there are any pieces between piece at
// at position p1 and position p2.
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

// doCastling goes through the steps to make sure that castling for king
// is legal and if it is, does the castling move, or returns an error
// explaining why it's not, if it isn't.
func (b *Board) doCastling(king *Piece, p1, p2 Pos) error {
	if b.check[king.Color] {
		return ErrCastleWithKingInCheck
	}

	if i, found := b.hasMoved[king]; found && i > 0 {
		return ErrKingOrRookMoved
	}

	switch p2.X {
	case 2: // Queen-side.
		// Move the queen-side rook to d1 or d8.
		piece, found := b.posToPiece[Pos{0, p2.Y}]
		if !found {
			return ErrNoRookToCastleWith
		}
		if i, found := b.hasMoved[piece]; found && i > 0 {
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

		// Add the rook to it's new position.
		b.posToPiece[Pos{3, p2.Y}] = piece

		// Remove the rook from the old position.
		delete(b.posToPiece, Pos{0, p2.Y})
	case 6: // King-side.
		// Move the king-side rook to f1 or f8.
		piece, found := b.posToPiece[Pos{7, p2.Y}]
		if !found {
			return ErrNoRookToCastleWith
		}
		if i, found := b.hasMoved[piece]; found && i > 0 {
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

		// Add the rook to it's new position.
		b.posToPiece[Pos{5, p2.Y}] = piece

		// Remove the rook from the old position.
		delete(b.posToPiece, Pos{7, p2.Y})
	default:
		// Shouldn't happen if called correctly.
		return fmt.Errorf("can't castle king to position %s", p2)
	}

	// Move the king to it's new position.
	//
	// The history will be able to tell that it was a castling
	// by which positions the king moved from and where to.
	b.makeMove(b.newMove(king, p1, p2, false))

	return nil
}
