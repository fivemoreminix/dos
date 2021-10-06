package dos

import "testing"

// Min returns the smaller of the two values.
func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// Max returns the larger of the two values.
func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// Clamp keeps the input value within a range of [min, max].
func Clamp(value, min, max int) int {
	return Max(min, Min(value, max))
}

type Rect struct {
	X, Y int
	W, H int
}

func (r Rect) HasPoint(x, y int) bool {
	return x >= r.X && y >= r.Y && x < r.X+r.W && y < r.Y+r.H
}

func TestRectHasPoint(t *testing.T) {
	if !(Rect{0, 0, 1, 1}.HasPoint(0, 0)) {
		t.Fail()
	}
	if (Rect{1, 1, 1, 1}.HasPoint(0, 0)) {
		t.Fail()
	}
	if (Rect{1, 1, 0, 0}.HasPoint(1, 1)) {
		t.Fail()
	}
	wide := Rect{3, 2, 15, 4}
	if wide.HasPoint(5, 1) {
		t.Fail()
	}
	if wide.HasPoint(8, 7) {
		t.Fail()
	}
	if !wide.HasPoint(5, 3) {
		t.Fail()
	}
}
