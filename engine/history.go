package engine

import "fmt"

func (b *Board) UndoMove() error {
	// Get the previous move from b.history.
	move, err := b.prevMove()
	if err != nil {
		return err
	}

	// Put the move's piece back to position from.
	b.posToPiece[move.From] = move.Piece

	// Delete the piece from position to.
	delete(b.posToPiece, move.To)

	// Put anything that was captured, back at position to.
	if move.Captured != nil {
		if move.EnPassant {
			b.posToPiece[Pos{move.To.X, move.From.Y}] = move.Captured
		} else {
			b.posToPiece[move.To] = move.Captured
		}
	}

	// Decrement b.hasMoved for move.Piece.
	b.hasMoved[move.Piece]--

	// If piece is a king, set it's position back to from.
	if move.Piece.Name == King {
		b.kings[move.Piece.Color] = move.From

		// Check if the move was a castling move.
		//
		// King side castling.
		if move.From.X == move.To.X+2 {
			// Put the rook back to Pos{0, move.To.Y}
			rook, found := b.posToPiece[Pos{move.To.X + 1, move.To.Y}]
			if !found {
				return fmt.Errorf("error: undo castling, rook not found")
			}
			b.posToPiece[Pos{0, move.To.Y}] = rook

			// Delete the rook from it's castled position.
			delete(b.posToPiece, Pos{move.To.X + 1, move.To.Y})
		}
		// Queen side castling.
		if move.From.X == move.To.X-2 {
			// Put the rook back to Pos{8, move.To.Y}
			rook, found := b.posToPiece[Pos{move.To.X - 1, move.To.Y}]
			if !found {
				return fmt.Errorf("error: undo castling, rook not found")
			}
			b.posToPiece[Pos{7, move.To.Y}] = rook

			// Delete the rook from it's castled position.
			delete(b.posToPiece, Pos{move.To.X - 1, move.To.Y})
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
	b.turn = move.Piece.Color

	// Decrement b.moveNum.
	b.moveNum--

	return nil
}

func (b *Board) prevMove() (*MoveInfo, error) {
	if b.moveNum < 0 {
		return nil, ErrNoPreviousMove
	}
	return b.history[b.moveNum], nil
}

// TODO: Redo move in history.
// func (b *Board) RedoMove() error {}
