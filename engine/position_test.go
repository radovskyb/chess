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
