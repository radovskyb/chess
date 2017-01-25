package engine

func (b *Board) UndoMove() error {
	if b.moveNum < 0 {
		return ErrHistoryIsEmpty
	}

	// Get the move from b.history at position b.moveNum.
	move := b.history[b.moveNum]

	// Put the move's piece back to position from.
	b.posToPiece[move.from] = move.piece

	// Delete the piece from position to.
	delete(b.posToPiece, move.to)

	// Put anything that was captured, back at position to.
	if move.captured != nil {
		if move.enPassant {
			b.posToPiece[Pos{move.to.X, move.from.Y}] = move.captured
		} else {
			b.posToPiece[move.to] = move.captured
		}
	}

	// Decrement b.hasMoved for move.piece.
	b.hasMoved[move.piece]--

	// If piece is a king, set it's position back to from.
	if move.piece.Name == King {
		b.kings[move.piece.Color] = move.from
	}

	// Set the turn to piece's color.
	b.turn = move.piece.Color

	// Decrement b.moveNum.
	b.moveNum--

	return nil
}

func (b *Board) prevMove() (*move, error) {
	if b.moveNum < 0 {
		return nil, ErrNoPreviousMove
	}
	return b.history[b.moveNum], nil
}

// TODO: Redo move in history.
// func (b *Board) RedoMove() error {}
