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

	for _, move := range setupMoves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	moveIntoCheckFrom, moveIntoCheckTo := "e1", "e2"
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

	for _, move := range setupMoves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	moveIntoCheckFrom, moveIntoCheckTo := "d2", "d3"
	err := b.MoveByLocation(moveIntoCheckFrom, moveIntoCheckTo)
	if err == nil {
		t.Errorf("expected to not be able to move from %s to %s",
			moveIntoCheckFrom, moveIntoCheckTo)
	}
}

func TestCanMovePieceAfterBlockingCheck(t *testing.T) {
	b := NewBoard()

	setupMoves := []struct {
		from, to string
	}{
		{"a2", "a3"},
		{"e7", "e5"},
		{"a3", "a4"},
		{"f8", "b4"},
		// Block king's line of sight with another piece.
		{"c2", "c3"},
		{"a7", "a5"},
	}

	for _, move := range setupMoves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	moveFrom, moveTo := "d2", "d3"
	err := b.MoveByLocation(moveFrom, moveTo)
	if err != nil {
		t.Errorf("expected to be able to move from %s to %s",
			moveFrom, moveTo)
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

func TestKingCanTakeCheckingPiece(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"e2", "e3"},
		{"e7", "e5"},
		{"b2", "b4"},
		{"f8", "b4"},
		{"a2", "a3"},
		{"b4", "d2"},
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

	if err := b.MoveByLocation("e1", "d2"); err != nil {
		t.Error("expected king to be able to take checking piece")
	}
}

func TestPieceCanTakeCheckingPiece(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"e2", "e3"},
		{"e7", "e5"},
		{"b2", "b4"},
		{"f8", "b4"},
		{"a2", "a3"},
		{"b4", "d2"},
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

	if err := b.MoveByLocation("c1", "d2"); err != nil {
		t.Error("expected piece to be able to take checking piece")
	}
}

func TestPositionAttacked(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"e2", "e3"},
		{"e7", "e5"},
		{"b2", "b4"},
		{"f8", "b4"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	attackedLocations := []string{
		"c3", "d2", "d4", "f4", "e7", "f6",
		"g5", "h4", "a3", "c5", "d6", "a5",
	}
	for _, loc := range attackedLocations {
		pos, err := locToPos(loc)
		if err != nil {
			t.Error(err)
		}
		if !b.positionAttacked(pos, Black) {
			t.Errorf("expected position %s to be attacked", pos)
		}
	}
}

func TestPositionAttackedFromStart(t *testing.T) {
	b := NewBoard()

	whiteAttackedLocations := []string{
		"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	}
	for _, loc := range whiteAttackedLocations {
		pos, err := locToPos(loc)
		if err != nil {
			t.Error(err)
		}
		if !b.positionAttacked(pos, White) {
			t.Errorf("expected position %s to be attacked", pos)
		}
	}

	blackAttackedLocations := []string{
		"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	}
	for _, loc := range blackAttackedLocations {
		pos, err := locToPos(loc)
		if err != nil {
			t.Error(err)
		}
		if !b.positionAttacked(pos, Black) {
			t.Errorf("expected position %s to be attacked", pos)
		}
	}

	notAttackedForBoth := []string{
		"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
		"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	}
	for _, loc := range notAttackedForBoth {
		pos, err := locToPos(loc)
		if err != nil {
			t.Error(err)
		}
		if b.positionAttacked(pos, Black) {
			t.Errorf("expected position %s to not be attacked", pos)
		}
		if b.positionAttacked(pos, White) {
			t.Errorf("expected position %s to not be attacked", pos)
		}
	}
}
