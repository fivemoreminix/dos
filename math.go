package dos

import "testing"

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
