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

func TestMoveByLocation(t *testing.T) {
	b := NewBoard()

	testCases := []struct {
		loc1, loc2 string
		name       PieceName
	}{
		{"a2", "a4", Pawn},
		{"a7", "a5", Pawn},
		{"a1", "a3", Rook},
		{"a8", "a6", Rook},
		{"b1", "c3", Knight},
		{"b8", "c6", Knight},
	}
	for _, tc := range testCases {
		if err := b.MoveByLocation(tc.loc1, tc.loc2); err != nil {
			t.Error(err)
		}
		piece, err := b.GetPieceAt(tc.loc2)
		if err != nil {
			t.Error(err)
		}
		if piece.Name != tc.name {
			t.Errorf("expected piece to be %v, got %v", tc.name, piece.Name)
		}
	}
}

func TestMove(t *testing.T) {
	b := NewBoard()

	testCases := []struct {
		p1, p2 Pos
		name   PieceName
	}{
		{Pos{0, 1}, Pos{0, 3}, Pawn},
		{Pos{0, 6}, Pos{0, 4}, Pawn},
		{Pos{0, 0}, Pos{0, 2}, Rook},
		{Pos{0, 7}, Pos{0, 5}, Rook},
		{Pos{1, 0}, Pos{2, 2}, Knight},
		{Pos{1, 7}, Pos{2, 5}, Knight},
	}
	for _, tc := range testCases {
		if err := b.Move(tc.p1, tc.p2); err != nil {
			t.Error(err)
		}
		piece, found := b.posToPiece[tc.p2]
		if !found {
			t.Errorf("expected to find piece %v at pos %v", tc.name, tc.p2)
		}
		if piece.Name != tc.name {
			t.Errorf("expected piece to be %v, got %v", tc.name, piece.Name)
		}
	}
}

func TestMoveLegal(t *testing.T) {
	b := NewBoard()

	testCases := []struct {
		p1, p2 Pos
		legal  bool
	}{
		// Test some white pieces.
		{Pos{0, 0}, Pos{0, 3}, false},
		{Pos{0, 1}, Pos{0, 3}, true},
		{Pos{1, 0}, Pos{1, 2}, false},
		{Pos{1, 0}, Pos{2, 2}, true},
		// Test some black pieces.
		{Pos{0, 7}, Pos{0, 6}, false},
		{Pos{0, 6}, Pos{0, 4}, true},
		{Pos{1, 7}, Pos{1, 5}, false},
		{Pos{1, 7}, Pos{2, 5}, true},
	}
	for _, tc := range testCases {
		piece, found := b.posToPiece[tc.p1]
		if !found {
			t.Errorf("no piece found at pos %v", tc.p1)
		}
		err := b.moveLegal(piece, tc.p1, tc.p2)
		if tc.legal && err != nil {
			t.Errorf("expected move from pos %v to pos %v be legal", tc.p1, tc.p2)
		}
		if !tc.legal && err == nil {
			t.Errorf("expected move from pos %v to pos %v be illegal", tc.p1, tc.p2)
		}
	}
}

func TestMoveBlocked(t *testing.T) {
	b := NewBoard()

	type pieceToPos struct {
		piece *Piece
		pos   Pos
	}

	testCases := []struct {
		setupMoves []pieceToPos
		blockTests [][2]Pos
	}{
		// White Pawn
		{
			[]pieceToPos{
				{&Piece{Pawn, White}, Pos{4, 3}},
				{&Piece{Pawn, White}, Pos{4, 4}}, // Pawn on top.
			},
			[][2]Pos{
				{Pos{4, 3}, Pos{4, 5}}, // Test blocking up
			},
		},
		// Black Pawn
		{
			[]pieceToPos{
				{&Piece{Pawn, Black}, Pos{4, 3}},
				{&Piece{Pawn, Black}, Pos{4, 2}}, // Pawn on bottom.
			},
			[][2]Pos{
				{Pos{4, 3}, Pos{4, 1}}, // Test blocking down
			},
		},
		// Rook
		{
			[]pieceToPos{
				{&Piece{Rook, White}, Pos{4, 3}},
				{&Piece{Pawn, White}, Pos{3, 3}}, // Pawn on left.
				{&Piece{Pawn, White}, Pos{5, 3}}, // Pawn on right.
				{&Piece{Pawn, White}, Pos{4, 4}}, // Pawn on top.
				{&Piece{Pawn, White}, Pos{4, 2}}, // Pawn on bottom.
			},
			[][2]Pos{
				{Pos{4, 3}, Pos{2, 3}}, // Test blocking left
				{Pos{4, 3}, Pos{6, 3}}, // Test blocking right
				{Pos{4, 3}, Pos{4, 6}}, // Test blocking up
				{Pos{4, 3}, Pos{4, 1}}, // Test blocking down
			},
		},
		// Bishop
		{
			[]pieceToPos{
				{&Piece{Bishop, White}, Pos{3, 3}},
				{&Piece{Pawn, White}, Pos{4, 4}}, // Pawn up+right.
				{&Piece{Pawn, White}, Pos{2, 4}}, // Pawn up+left.
				{&Piece{Pawn, White}, Pos{4, 2}}, // Pawn down+right.
				{&Piece{Pawn, White}, Pos{2, 2}}, // Pawn down+left.
			},
			[][2]Pos{
				{Pos{3, 3}, Pos{5, 5}}, // Test blocking up+right.
				{Pos{3, 3}, Pos{1, 5}}, // Test blocking up+left.
				{Pos{3, 3}, Pos{5, 1}}, // Test blocking down+right.
				{Pos{3, 3}, Pos{1, 1}}, // Test blocking down+left.
			},
		},
		// Queen
		{
			[]pieceToPos{
				{&Piece{Queen, White}, Pos{3, 3}},
				// Diagonal blockages.
				{&Piece{Pawn, White}, Pos{4, 4}}, // Pawn up+right.
				{&Piece{Pawn, White}, Pos{2, 4}}, // Pawn up+left.
				{&Piece{Pawn, White}, Pos{4, 2}}, // Pawn down+right.
				{&Piece{Pawn, White}, Pos{2, 2}}, // Pawn down+left.
				// Line blockages.
				{&Piece{Pawn, White}, Pos{4, 3}}, // Pawn right.
				{&Piece{Pawn, White}, Pos{2, 3}}, // Pawn left.
				{&Piece{Pawn, White}, Pos{3, 4}}, // Pawn up.
				{&Piece{Pawn, White}, Pos{3, 2}}, // Pawn down.
			},
			[][2]Pos{
				// Diagonal blockages.
				{Pos{3, 3}, Pos{5, 5}}, // Test blocking up+right.
				{Pos{3, 3}, Pos{1, 5}}, // Test blocking up+left.
				{Pos{3, 3}, Pos{5, 1}}, // Test blocking down+right.
				{Pos{3, 3}, Pos{1, 1}}, // Test blocking down+left.
				// Line blockages.
				{Pos{3, 3}, Pos{5, 3}}, // Test blocking right.
				{Pos{3, 3}, Pos{1, 3}}, // Test blocking left.
				{Pos{3, 3}, Pos{3, 5}}, // Test blocking up.
				{Pos{3, 3}, Pos{3, 1}}, // Test blocking down.
			},
		},
	}
	for _, tc := range testCases {
		b.clear() // clear the board

		for _, move := range tc.setupMoves {
			b.posToPiece[move.pos] = move.piece
		}
		for _, positions := range tc.blockTests {
			piece, found := b.posToPiece[positions[0]]
			if !found {
				t.Errorf("no piece found at pos %v", positions[0])
			}
			if !b.moveBlocked(piece, positions[0], positions[1]) {
				t.Errorf("expected moving %s from %v to %v to block",
					piece, positions[0], positions[1])
			}
		}
	}
}

// func TestMoveNotBlocked(t *testing.T) {}
