package dos

import (
	"github.com/gdamore/tcell/v2"
)

// A Row orders its children horizontally.
type Row struct {
	Children      []Widget
	VerticalAlign Alignment
	FocusedIndex  int // Index of child that receives focus
	OnKeyEvent    func(ev *tcell.EventKey) bool
}

func (r *Row) GetChildRects(rect Rect) []Rect {
	if childLen := len(r.Children); childLen > 0 {
		rects := make([]Rect, childLen)
		childWidthSum := 0
		childMaxWidth := rect.W / childLen
		for i := 0; i < childLen; i++ {
			w, h := r.Children[i].DisplaySize(childMaxWidth, rect.H)
			// TODO: implement VerticalAlign for Row here
			rects[i] = Rect{0, rect.Y, w, h}
			childWidthSum += w
		}
		if childWidthSum < rect.H {
			qualifyingChildren := 0
			for i := 0; i < childLen; i++ {
				if rects[i].W < childMaxWidth {
					qualifyingChildren++
				}
			}
			growAmount := (rect.W - childWidthSum) / qualifyingChildren
			for i := 0; i < childLen; i++ {
				if rects[i].W == childMaxWidth {
					rects[i].W += growAmount
				}
			}
		}
		// Stack rects one on top of the other
		childWidthSum = 0
		for i := 0; i < childLen; i++ {
			rects[i].X = rect.X + childWidthSum
			childWidthSum += rects[i].W
		}
		return rects
	}
	return nil
}

func (r *Row) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	rects := r.GetChildRects(currentRect)
	for i := range r.Children {
		if r.Children[i].HandleMouse(rects[i], ev) {
			r.SetFocused(false) // Unfocus any prior-focused child
			r.FocusedIndex = i
			return true
		}
	}
	return false
}

func (r *Row) HandleKey(ev *tcell.EventKey) bool {
	if r.OnKeyEvent != nil && r.OnKeyEvent(ev) {
		return true
	}
	for i := range r.Children {
		if r.Children[i].HandleKey(ev) {
			return true
		}
	}
	return false
}

func (r *Row) SetFocused(b bool) {
	if r.FocusedIndex < len(r.Children) {
		r.Children[r.FocusedIndex].SetFocused(b)
	}
}

func (r *Row) DisplaySize(boundsW, boundsH int) (w, h int) {
	rects := r.GetChildRects(Rect{0, 0, boundsW, boundsH})
	if rects == nil {
		return 0, 0
	}
	height := 0
	width := 0
	for i := range rects {
		width += rects[i].W // combined width
		if rects[i].H > height {
			height = rects[i].H // only the maximum height
		}
	}
	return width, height
}

func (r *Row) Draw(rect Rect, s tcell.Screen) {
	rects := r.GetChildRects(rect)
	for i := range rects {
		r.Children[i].Draw(rects[i], s)
	}
}
