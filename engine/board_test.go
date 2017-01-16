package engine

import "testing"

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	// The initial turn should be white.
	if b.Turn != White {
		t.Error("expected turn to be white")
	}

	// Check if there are the correct amount of pieces on the new board.
	if numPieces := len(b.posToPiece); numPieces != 32 {
		t.Errorf("expected len(posToPiece) to be 32, got %d", numPieces)
	}

	// Make sure that starting positions match the correct piece name.
	posToPieceNames := []struct {
		pos  Pos
		name PieceName
	}{
		// Test some white pieces.
		{Pos{0, 1}, Pawn},
		{Pos{0, 0}, Rook},
		{Pos{1, 0}, Knight},
		{Pos{2, 0}, Bishop},
		{Pos{3, 0}, Queen},
		{Pos{4, 0}, King},
		// Test some black pieces.
		{Pos{0, 6}, Pawn},
		{Pos{0, 7}, Rook},
		{Pos{1, 7}, Knight},
		{Pos{2, 7}, Bishop},
		{Pos{3, 7}, Queen},
		{Pos{4, 7}, King},
	}
	for _, tc := range posToPieceNames {
		piece, found := b.posToPiece[tc.pos]
		if !found {
			t.Errorf("piece not found at position %v", tc.pos)
		}
		if piece.Name != tc.name {
			t.Errorf("expected piece to be %v, got %v", tc.name, piece.Name)
		}
	}
}

func TestGetPieceAt(t *testing.T) {
	b := NewBoard()

	testCases := []struct {
		loc  string
		name PieceName
	}{
		// Test some white pieces.
		{"a2", Pawn},
		{"a1", Rook},
		{"b1", Knight},
		{"c1", Bishop},
		{"d1", Queen},
		{"e1", King},
		// Test some black pieces.
		{"a7", Pawn},
		{"a8", Rook},
		{"b8", Knight},
		{"c8", Bishop},
		{"d8", Queen},
		{"e8", King},
	}
	for _, tc := range testCases {
		piece, err := b.GetPieceAt(tc.loc)
		if err != nil {
			t.Error(err)
		}
		if piece.Name != tc.name {
			t.Errorf("expected piece to be %v, got %v", tc.name, piece.Name)
		}
	}
}
