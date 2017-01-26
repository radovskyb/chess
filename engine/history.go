package engine

import "fmt"

func (b *Board) UndoMove() error {
	// Get the previous move from b.history.
	move, err := b.prevMove()
	if err != nil {
		return err
	}

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

		// Check if the move was a castling move.
		//
		// King side castling.
		if move.from.X == move.to.X+2 {
			// Put the rook back to Pos{0, move.to.Y}
			rook, found := b.posToPiece[Pos{move.to.X + 1, move.to.Y}]
			if !found {
				return fmt.Errorf("error: undo castling, rook not found")
			}
			b.posToPiece[Pos{0, move.to.Y}] = rook

			// Delete the rook from it's castled position.
			delete(b.posToPiece, Pos{move.to.X + 1, move.to.Y})
		}
		// Queen side castling.
		if move.from.X == move.to.X-2 {
			// Put the rook back to Pos{8, move.to.Y}
			rook, found := b.posToPiece[Pos{move.to.X - 1, move.to.Y}]
			if !found {
				return fmt.Errorf("error: undo castling, rook not found")
			}
			b.posToPiece[Pos{7, move.to.Y}] = rook

			// Delete the rook from it's castled position.
			delete(b.posToPiece, Pos{move.to.X - 1, move.to.Y})
		}
	}

	// Set the checks on the board back to the previous move's checks.
	b.check[White], b.check[Black] = false, false
	if b.kingInCheck(White) {
		b.check[White] = true
	}
	if b.kingInCheck(Black) {
		b.check[Black] = true
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
