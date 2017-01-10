package engine

import "testing"

func TestLocToPos(t *testing.T) {
	testCases := []struct {
		loc string
		pos Pos
	}{
		// Test some low end locations.
		{"a1", Pos{0, 0}},
		{"b1", Pos{1, 0}},
		{"c1", Pos{2, 0}},
		{"a2", Pos{0, 1}},
		{"b2", Pos{1, 1}},
		{"c2", Pos{2, 1}},
		// Test some high end locations.
		{"h8", Pos{7, 7}},
		{"g8", Pos{6, 7}},
		{"f8", Pos{5, 7}},
		{"h7", Pos{7, 6}},
		{"g7", Pos{6, 6}},
		{"f7", Pos{5, 6}},
	}
	for _, tc := range testCases {
		pos, err := locToPos(tc.loc)
		if err != nil {
			t.Error(err)
		}
		if pos != tc.pos {
			t.Errorf("expected pos to be %v, got %v", tc.pos, pos)
		}
	}
}

func TestDiagPositions(t *testing.T) {
	testCases := []struct {
		cur      Pos
		expected []Pos
	}{
		{
			Pos{3, 3}, []Pos{
				{0, 0}, {1, 1}, {2, 2},
				{4, 4}, {5, 5}, {6, 6},
				{2, 4}, {1, 5}, {0, 6},
				{6, 0}, {3, 3}, {7, 7},
				{2, 2}, {5, 1},
			},
		},
		{
			Pos{5, 6}, []Pos{
				{0, 1}, {1, 2}, {2, 3},
				{3, 4}, {4, 5}, {5, 6},
				{6, 5}, {6, 7}, {4, 7},
				{7, 4},
			},
		},
		{
			Pos{1, 1}, []Pos{
				{0, 0}, {1, 1}, {2, 2},
				{3, 3}, {4, 4}, {5, 5},
				{6, 6}, {7, 7}, {2, 0},
				{0, 2},
			},
		},
	}
	for _, tc := range testCases {
		positions := diagPositions(tc.cur)
		if len(positions) != len(tc.expected) {
			t.Errorf("expected len(positions) to be %d", len(tc.expected))
		}
		for _, pos := range tc.expected {
			if _, found := positions[pos]; !found {
				t.Errorf("expected to find pos %v", pos)
			}
		}
	}
}

func TestLinePositions(t *testing.T) {
	testCases := []struct {
		cur      Pos
		expected []Pos
	}{
		{
			Pos{3, 3}, []Pos{
				// Veritcal line positions.
				{3, 0}, {3, 1}, {3, 2}, {3, 3},
				{3, 4}, {3, 5}, {3, 6}, {3, 7},
				// Horizontal line positions.
				{0, 3}, {1, 3}, {2, 3}, {4, 3},
				{5, 3}, {6, 3}, {7, 3},
			},
		},
		{
			Pos{5, 1}, []Pos{
				// Veritcal line positions.
				{5, 0}, {5, 1}, {5, 2}, {5, 3},
				{5, 4}, {5, 5}, {5, 6}, {5, 7},
				// Horizontal line positions.
				{0, 1}, {1, 1}, {2, 1}, {3, 1},
				{4, 1}, {6, 1}, {7, 1},
			},
		},
		{
			Pos{6, 2}, []Pos{
				// Veritcal line positions.
				{6, 0}, {6, 1}, {6, 2}, {6, 3},
				{6, 4}, {6, 5}, {6, 6}, {6, 7},
				// Horizontal line positions.
				{0, 2}, {1, 2}, {2, 2}, {3, 2},
				{4, 2}, {5, 2}, {7, 2},
			},
		},
	}
	for _, tc := range testCases {
		positions := linePositions(tc.cur)
		if len(positions) != len(tc.expected) {
			t.Errorf("expected len(positions) to be %d", len(tc.expected))
		}
		for _, pos := range tc.expected {
			if _, found := positions[pos]; !found {
				t.Errorf("expected to find pos %v", pos)
			}
		}
	}
}
