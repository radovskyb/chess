package engine

import "testing"

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

func TestHasCheck(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"e2", "e3"},
		{"f7", "f5"},
		{"a2", "a4"},
		{"e8", "f7"},
		{"d1", "e2"},
		{"a7", "a5"},
		{"e2", "c4"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	hasCheck, color := b.HasCheck()
	if !hasCheck {
		t.Error("expected board to have a check")
	}
	if color != Black {
		t.Error("expected black to be in check")
	}
}

func TestCantMoveKingIntoCheck(t *testing.T) {
	b := NewBoard()

	setupMoves := []struct {
		from, to string
	}{
		{"e2", "e4"},
		{"d7", "d5"},
		{"a2", "a3"},
		{"c8", "g4"},
	}

	moveIntoCheckFrom, moveIntoCheckTo := "e1", "e2"

	for _, move := range setupMoves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	err := b.MoveByLocation(moveIntoCheckFrom, moveIntoCheckTo)
	if err == nil {
		t.Errorf("expected to not be able to move from %s to %s",
			moveIntoCheckFrom, moveIntoCheckTo)
	}
}

func TestCantMovePieceIntoCheck(t *testing.T) {
	b := NewBoard()

	setupMoves := []struct {
		from, to string
	}{
		{"a2", "a3"},
		{"e7", "e5"},
		{"a3", "a4"},
		{"f8", "b4"},
	}

	moveIntoCheckFrom, moveIntoCheckTo := "d2", "d3"

	for _, move := range setupMoves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	err := b.MoveByLocation(moveIntoCheckFrom, moveIntoCheckTo)
	if err == nil {
		t.Errorf("expected to not be able to move from %s to %s",
			moveIntoCheckFrom, moveIntoCheckTo)
	}
}

func TestCantMovePieceWhenInCheck(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"e2", "e3"},
		{"f7", "f5"},
		{"a2", "a4"},
		{"e8", "f7"},
		{"d1", "e2"},
		{"a7", "a5"},
		{"e2", "c4"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	hasCheck, _ := b.HasCheck()
	if !hasCheck {
		t.Error("expected board to have a check")
	}

	err := b.MoveByLocation("a8", "a6")
	if err == nil {
		t.Error("expected to not be able to move piece whilst in check")
	}
}

func TestCanMovePieceInCheckToUncheck(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"e2", "e3"},
		{"f7", "f5"},
		{"a2", "a4"},
		{"e8", "f7"},
		{"d1", "e2"},
		{"a7", "a5"},
		{"e2", "c4"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	hasCheck, _ := b.HasCheck()
	if !hasCheck {
		t.Error("expected board to have a check")
	}

	if err := b.MoveByLocation("e7", "e6"); err != nil {
		t.Error("expected to be able to move piece to block check")
	}
}
