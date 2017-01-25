package engine

import "testing"

func TestHistoryEmpty(t *testing.T) {
	b := NewBoard()

	// No moves have been made, expect an ErrNoPreviousMove error.
	if err := b.UndoMove(); err == nil {
		t.Error("expected error to be ErrNoPreviousMove")
	}

	// Make a move.
	if err := b.MoveByLocation("a2", "a3"); err != nil {
		t.Fatalf("moving from a2 to a3 failed: %s", err.Error())
	}

	// Undo the previous move.
	if err := b.UndoMove(); err != nil {
		t.Error("expected undo to work: %s", err.Error())
	}

	// Undo a second time, which this time, shouldn't work.
	if err := b.UndoMove(); err == nil {
		t.Error("expected error to be ErrNoPreviousMove")
	}
}

func TestUndoNormalMove(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"d2", "d3"},
		{"d7", "d6"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	// Make sure there's a pawn at location d6.
	piece, err := b.GetPieceAt("d6")
	if err != nil || piece.Name != Pawn {
		t.Error("expected there to be a pawn at location d6")
	}

	if err := b.UndoMove(); err != nil {
		t.Error("expected to be able to undo previous move")
	}

	// Make sure there's no longer a pawn at location d6.
	piece, err = b.GetPieceAt("d6")
	if err == nil {
		t.Error("expected there to be no pawn at location d6")
	}
}

func TestUndoTakePiece(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"d2", "d4"},
		{"e7", "e5"},
		{"d4", "e5"}, // Take the pawn.
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	// Check that a piece was taken.
	if len(b.posToPiece) != 31 {
		t.Errorf("expected there to only be 31 pieces, got %d",
			len(b.posToPiece))
	}

	// Make sure there's no longer a pawn at location d4.
	_, err := b.GetPieceAt("d4")
	if err == nil {
		t.Error("expected there to be no pawn at location d4")
	}

	if err := b.UndoMove(); err != nil {
		t.Error("expected to be able to undo previous move")
	}

	// Check that there's 32 pieces again.
	if len(b.posToPiece) != 32 {
		t.Errorf("expected there to be 32 pieces, got %d",
			len(b.posToPiece))
	}

	// Make sure there's a pawn at d4.
	piece, err := b.GetPieceAt("d4")
	if err != nil || piece.Name != Pawn {
		t.Error("expected there to be a pawn at location d4")
	}

	// Make sure there's a pawn at e5.
	piece, err = b.GetPieceAt("e5")
	if err != nil || piece.Name != Pawn {
		t.Error("expected there to be a pawn at location e5")
	}
}

func TestUndoCastlingQueenSide(t *testing.T) {
	b := NewBoard()
	b.clear()

	// White pieces.
	wK := &Piece{King, White}
	wR := &Piece{Rook, White}

	// Add king and rook on the queen-side for white.
	b.posToPiece[Pos{4, 0}] = wK
	b.posToPiece[Pos{0, 0}] = wR

	// Castle white by moving the king over 2 squares queen-side.
	if err := b.MoveByLocation("e1", "c1"); err != nil {
		t.Error(err)
	}

	// Check that the white king is in the right place.
	piece, err := b.GetPieceAt("c1")
	if err != nil {
		t.Error(err)
	}
	if piece != wK {
		t.Errorf("expected piece to be white king, got %s", piece)
	}
	// Check that the white rook is in the right place.
	piece, err = b.GetPieceAt("d1")
	if err != nil {
		t.Error(err)
	}
	if piece != wR {
		t.Errorf("expected piece to be white rook, got %s", piece)
	}

	// Undo the castling.
	if err := b.UndoMove(); err != nil {
		t.Error(err)
	}

	// Make sure there's no longer a king at location c1.
	_, err = b.GetPieceAt("c1")
	if err == nil {
		t.Error("expected there to be no piece at location c1")
	}
	// Make sure there's no longer a rook at location d1.
	_, err = b.GetPieceAt("d1")
	if err == nil {
		t.Error("expected there to be no piece at location d1")
	}

	// Make sure the king is back at location e1.
	piece, err = b.GetPieceAt("e1")
	if err != nil || piece.Name != King {
		t.Error("expected there to be a king at location e1")
	}
	// Make sure the rook is back at location a1.
	piece, err = b.GetPieceAt("a1")
	if err != nil || piece.Name != Rook {
		t.Error("expected there to be a rook at location a1")
	}
}

func TestUndoCastlingKingSide(t *testing.T) {
	b := NewBoard()
	b.clear()

	// Switch to black's turn.
	b.turn ^= 1

	// Black pieces.
	bK := &Piece{King, Black}
	bR := &Piece{Rook, Black}

	// Add king and rook on the king-side for black.
	b.posToPiece[Pos{4, 7}] = bK
	b.posToPiece[Pos{7, 7}] = bR

	// Castle black by moving the king over 2 squares king-side.
	if err := b.MoveByLocation("e8", "g8"); err != nil {
		t.Error(err)
	}

	// Check that the black king is in the right place.
	piece, err := b.GetPieceAt("g8")
	if err != nil {
		t.Error(err)
	}
	if piece != bK {
		t.Errorf("expected piece to be black king, got %s", piece)
	}
	// Check that the black rook is in the right place.
	piece, err = b.GetPieceAt("f8")
	if err != nil {
		t.Error(err)
	}
	if piece != bR {
		t.Errorf("expected piece to be black rook, got %s", piece)
	}

	// Undo the castling.
	if err := b.UndoMove(); err != nil {
		t.Error(err)
	}

	// Make sure there's no longer a king at location g8.
	_, err = b.GetPieceAt("g8")
	if err == nil {
		t.Error("expected there to be no piece at location g8")
	}
	// Make sure there's no longer a rook at location f8.
	_, err = b.GetPieceAt("f8")
	if err == nil {
		t.Error("expected there to be no piece at location f8")
	}

	// Make sure the king is back at location e8.
	piece, err = b.GetPieceAt("e8")
	if err != nil || piece.Name != King {
		t.Error("expected there to be a king at location e8")
	}
	// Make sure the rook is back at location h8.
	piece, err = b.GetPieceAt("h8")
	if err != nil || piece.Name != Rook {
		t.Error("expected there to be a rook at location h8")
	}
}

// TODO: Undo en passant.
// TODO: Undo promotion.
//
// TODO: Test prevMove
