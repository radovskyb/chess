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

func TestMoveBlocked(t *testing.T) {
	b := NewBoard()

	testCases := []struct {
		setupMoves []piecePos
		blockTests [][2]Pos
	}{
		// White Pawn
		{
			[]piecePos{
				{&Piece{Pawn, White}, Pos{4, 3}},
				{&Piece{Pawn, White}, Pos{4, 4}}, // Pawn on top.
			},
			[][2]Pos{
				{Pos{4, 3}, Pos{4, 5}}, // Test blocking up
			},
		},
		// Black Pawn
		{
			[]piecePos{
				{&Piece{Pawn, Black}, Pos{4, 3}},
				{&Piece{Pawn, Black}, Pos{4, 2}}, // Pawn on bottom.
			},
			[][2]Pos{
				{Pos{4, 3}, Pos{4, 1}}, // Test blocking down
			},
		},
		// Rook
		{
			[]piecePos{
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
			[]piecePos{
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
			[]piecePos{
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
			b.posToPiece[move.Pos] = move.Piece
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

// TODO: func TestMoveNotBlocked(t *testing.T) {}

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

func TestCastlingQueenSide(t *testing.T) {
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
}

func TestCastlingKingSide(t *testing.T) {
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
}

func TestCantCastleIfKingMoved(t *testing.T) {
	b := NewBoard()
	b.clear()

	// White pieces.
	wK := &Piece{King, White}
	wR := &Piece{Rook, White}

	// Black pieces.
	bK := &Piece{King, Black}
	bR := &Piece{Rook, Black}

	// Add king and rook on the quuen-side for white.
	b.posToPiece[Pos{4, 0}] = wK
	b.posToPiece[Pos{0, 0}] = wR

	// Add king and rook on the queen-side for black.
	b.posToPiece[Pos{4, 7}] = bK
	b.posToPiece[Pos{0, 7}] = bR

	// Move the white king.
	if err := b.MoveByLocation("e1", "d1"); err != nil {
		t.Error(err)
	}

	// Move the black king.
	if err := b.MoveByLocation("e8", "d8"); err != nil {
		t.Error(err)
	}

	// Move the white king back.
	if err := b.MoveByLocation("d1", "e1"); err != nil {
		t.Error(err)
	}

	// Move the black king back.
	if err := b.MoveByLocation("d8", "e8"); err != nil {
		t.Error(err)
	}

	// Try to castle white.
	if err := b.MoveByLocation("e1", "c1"); err != ErrKingOrRookMoved {
		t.Error("shouldn't be able to castle white, king has moved before")
	}

	// Switch to black's turn.
	b.turn ^= 1

	// Try to castle black.
	if err := b.MoveByLocation("e8", "c8"); err != ErrKingOrRookMoved {
		t.Error("shouldn't be able to castle black, king has moved before")
	}
}

func TestCantCastleIfRookMoved(t *testing.T) {
	b := NewBoard()
	b.clear()

	// White pieces.
	wK := &Piece{King, White}
	wR := &Piece{Rook, White}

	// Black pieces.
	bK := &Piece{King, Black}
	bR := &Piece{Rook, Black}

	// Add king and rook on the quuen-side for white.
	b.posToPiece[Pos{4, 0}] = wK
	b.posToPiece[Pos{0, 0}] = wR

	// Add king and rook on the queen-side for black.
	b.posToPiece[Pos{4, 7}] = bK
	b.posToPiece[Pos{0, 7}] = bR

	// Move the white rook.
	if err := b.MoveByLocation("a1", "b1"); err != nil {
		t.Error(err)
	}

	// Move the black rook.
	if err := b.MoveByLocation("a8", "b8"); err != nil {
		t.Error(err)
	}

	// Move the white king back.
	if err := b.MoveByLocation("b1", "a1"); err != nil {
		t.Error(err)
	}

	// Move the black king back.
	if err := b.MoveByLocation("b8", "a8"); err != nil {
		t.Error(err)
	}

	// Try to castle white.
	if err := b.MoveByLocation("e1", "c1"); err != ErrKingOrRookMoved {
		t.Error("shouldn't be able to castle white, king has moved before")
	}

	// Switch to black's turn.
	b.turn ^= 1

	// Try to castle black.
	if err := b.MoveByLocation("e8", "c8"); err != ErrKingOrRookMoved {
		t.Error("shouldn't be able to castle black, king has moved before")
	}
}

func TestCantCastleWhiteIfKingInCheck(t *testing.T) {
	b := NewBoard()
	b.clear()

	// Add king and rook on the quuen-side for white.
	b.posToPiece[Pos{4, 0}] = &Piece{King, White}
	b.posToPiece[Pos{0, 0}] = &Piece{Rook, White}

	b.turn ^= 1

	// Put black rook on board and move it to put white king in check.
	b.posToPiece[Pos{4, 4}] = &Piece{Rook, Black}
	if err := b.Move(Pos{4, 4}, Pos{4, 3}); err != nil {
		t.Error(err)
	}

	// Try to castle white.
	if err := b.MoveByLocation("e1", "c1"); err != ErrCastleWithKingInCheck {
		t.Error("shouldn't be able to castle white, king is in check")
	}
}

func TestCantCastleBlackIfKingInCheck(t *testing.T) {
	b := NewBoard()
	b.clear()

	// Add king and rook on the queen-side for black.
	b.posToPiece[Pos{4, 7}] = &Piece{King, Black}
	b.posToPiece[Pos{0, 7}] = &Piece{Rook, Black}

	// Put white rook on board and move it to put black king in check.
	b.posToPiece[Pos{4, 4}] = &Piece{Rook, White}
	if err := b.Move(Pos{4, 4}, Pos{4, 3}); err != nil {
		t.Error(err)
	}

	// Try to castle black.
	if err := b.MoveByLocation("e8", "c8"); err != ErrCastleWithKingInCheck {
		t.Error("shouldn't be able to castle black, king is in check")
	}
}

func TestCantCastleIfPieceBetweenKingAndRook(t *testing.T) {
	b := NewBoard()
	b.clear()

	// Add king and rook on the queen-side for white.
	b.posToPiece[Pos{4, 0}] = &Piece{King, White}
	b.posToPiece[Pos{0, 0}] = &Piece{Rook, White}
	b.posToPiece[Pos{7, 0}] = &Piece{Rook, White}

	// Put a bishop between the king and the rook on both sides.
	b.posToPiece[Pos{2, 0}] = &Piece{Bishop, White}
	b.posToPiece[Pos{6, 0}] = &Piece{Bishop, White}

	// Try to castle white queen-side.
	if err := b.MoveByLocation("e1", "c1"); err != ErrCastleWithPieceBetween {
		t.Error("shouldn't be able to castle white with a piece in the way")
	}

	// Try to castle white king-side.
	if err := b.MoveByLocation("e1", "g1"); err != ErrCastleWithPieceBetween {
		t.Error("shouldn't be able to castle white with a piece in the way")
	}

	b.turn ^= 1

	// Add king and rook on the queen-side for black.
	b.posToPiece[Pos{4, 7}] = &Piece{King, Black}
	b.posToPiece[Pos{0, 7}] = &Piece{Rook, Black}
	b.posToPiece[Pos{7, 7}] = &Piece{Rook, Black}

	// Put a bishop between the king and the rook on both sides.
	b.posToPiece[Pos{2, 7}] = &Piece{Bishop, Black}
	b.posToPiece[Pos{6, 7}] = &Piece{Bishop, Black}

	// Try to castle black queen-side.
	if err := b.MoveByLocation("e8", "c8"); err != ErrCastleWithPieceBetween {
		t.Error("shouldn't be able to castle black with a piece in the way")
	}

	// Try to castle black king-side.
	if err := b.MoveByLocation("e8", "g8"); err != ErrCastleWithPieceBetween {
		t.Error("shouldn't be able to castle black with a piece in the way")
	}
}

func TestCantCastleIfKingMovesThroughCheck(t *testing.T) {
	b := NewBoard()
	b.clear()

	// Add king and rook on the queen-side for white.
	b.posToPiece[Pos{4, 0}] = &Piece{King, White}
	b.posToPiece[Pos{0, 0}] = &Piece{Rook, White}
	b.posToPiece[Pos{7, 0}] = &Piece{Rook, White}

	// Put a bishop that will make both sides castle through check.
	b.posToPiece[Pos{4, 2}] = &Piece{Bishop, Black}

	// Try to castle white queen-side.
	if err := b.MoveByLocation("e1", "c1"); err != ErrCastleMoveThroughCheck {
		t.Error("shouldn't be able to castle when king moves through check")
	}

	// Try to castle white king-side.
	if err := b.MoveByLocation("e1", "g1"); err != ErrCastleMoveThroughCheck {
		t.Error("shouldn't be able to castle when king moves through check")
	}

	b.turn ^= 1

	// Add king and rook on the queen-side for black.
	b.posToPiece[Pos{4, 7}] = &Piece{King, Black}
	b.posToPiece[Pos{0, 7}] = &Piece{Rook, Black}
	b.posToPiece[Pos{7, 7}] = &Piece{Rook, Black}

	// Put a bishop that will make both sides castle through check.
	b.posToPiece[Pos{4, 5}] = &Piece{Bishop, White}

	// Try to castle black queen-side.
	if err := b.MoveByLocation("e8", "c8"); err != ErrCastleMoveThroughCheck {
		t.Error("shouldn't be able to castle when king moves through check")
	}

	// Try to castle black king-side.
	if err := b.MoveByLocation("e8", "g8"); err != ErrCastleMoveThroughCheck {
		t.Error("shouldn't be able to castle when king moves through check")
	}
}

func TestEnPassant(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"a2", "a4"},
		{"a7", "a6"},
		{"a4", "a5"},
		{"b7", "b5"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	// Make sure pawn is at b5.
	piece, err := b.GetPieceAt("b5")
	if err != nil {
		t.Error(err)
	}
	if piece.Name != Pawn {
		t.Error("expected to find pawn at location b5")
	}

	enPassantFrom, enPassantTo := "a5", "b6"
	if err := b.MoveByLocation(enPassantFrom, enPassantTo); err != nil {
		t.Errorf("expected en passant from %s to %s to work",
			enPassantFrom, enPassantTo)
	}

	// Make sure pawn at b5 was taken.
	piece, err = b.GetPieceAt("b5")
	if err == nil {
		t.Error("expected there to be no pawn at location b5")
	}
}

func TestCantEnPassantIfNotCorrectPosition(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"c2", "c4"},
		{"d7", "d5"},
		{"c4", "c5"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	enPassantFrom, enPassantTo := "d5", "c4"
	if err := b.MoveByLocation(enPassantFrom, enPassantTo); err == nil {
		t.Errorf("expected en passant from %s to %s not to work",
			enPassantFrom, enPassantTo)
	}

	if err := b.MoveByLocation("a7", "a6"); err != nil {
		t.Fatalf("moving from a7 to a6 failed: %s", err.Error())
	}

	enPassantFrom, enPassantTo = "c5", "d6"
	if err := b.MoveByLocation(enPassantFrom, enPassantTo); err == nil {
		t.Errorf("expected en passant from %s to %s not to work",
			enPassantFrom, enPassantTo)
	}
}

func TestCantEnPassantIfCausesOwnCheck(t *testing.T) {
	b := NewBoard()

	moves := []struct {
		from, to string
	}{
		{"d2", "d3"},
		{"d7", "d6"},
		{"c2", "c4"},
		{"d8", "d7"},
		{"c4", "c5"},
		{"d6", "c5"},
		{"e1", "d2"},
		{"a7", "a5"},
		{"d3", "d4"},
		{"a5", "a4"},
		{"d4", "d5"},
		{"e7", "e5"},
	}

	for _, move := range moves {
		if err := b.MoveByLocation(move.from, move.to); err != nil {
			t.Fatalf("moving from %s to %s failed: %s",
				move.from, move.to, err.Error())
		}
	}

	enPassantFrom, enPassantTo := "d5", "e6" // Causes own king to be in check.
	if err := b.MoveByLocation(enPassantFrom, enPassantTo); err == nil {
		t.Errorf("expected en passant from %s to %s not to work",
			enPassantFrom, enPassantTo)
	}
}
