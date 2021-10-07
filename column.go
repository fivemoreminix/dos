package dos

import (
	"github.com/gdamore/tcell/v2"
)

// A Column orders its children vertically.
type Column struct {
	Children        []Widget
	HorizontalAlign Alignment
	FocusedIndex    int // Index of child that receives focus
	OnKeyEvent      func(ev *tcell.EventKey) bool
}

func (c *Column) GetChildRects(rect Rect) []Rect {
	if childLen := len(c.Children); childLen > 0 {
		rects := make([]Rect, childLen)
		childHeightSum := 0
		childMaxHeight := rect.H / childLen
		for i := 0; i < childLen; i++ {
			w, h := c.Children[i].DisplaySize(rect.W, childMaxHeight)
			// TODO: implement HorizontalAlign for Column here
			rects[i] = Rect{rect.X, 0, w, h}
			childHeightSum += h
		}
		if childHeightSum < rect.H {
			qualifyingChildren := 0
			for i := 0; i < childLen; i++ {
				if rects[i].H < childMaxHeight {
					qualifyingChildren++
				}
			}
			if qualifyingChildren != 0 {
				growAmount := (rect.H - childHeightSum) / qualifyingChildren
				for i := 0; i < childLen; i++ {
					if rects[i].H == childMaxHeight {
						rects[i].H += growAmount
					}
				}
			}
		}
		// Stack rects one on top of the other
		childHeightSum = 0
		for i := 0; i < childLen; i++ {
			rects[i].Y = rect.Y + childHeightSum
			childHeightSum += rects[i].H
		}
		return rects
	}
	return nil
}

func (c *Column) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	rects := c.GetChildRects(currentRect)
	for i := range c.Children {
		if c.Children[i].HandleMouse(rects[i], ev) {
			c.SetFocused(false) // Unfocus any prior-focused child
			c.FocusedIndex = i
			return true
		}
	}
	return false
}

func (c *Column) HandleKey(ev *tcell.EventKey) bool {
	if c.OnKeyEvent != nil && c.OnKeyEvent(ev) {
		return true
	}
	for i := range c.Children {
		if c.Children[i].HandleKey(ev) {
			return true
		}
	}
	return false
}

func (c *Column) SetFocused(b bool) {
	if c.FocusedIndex < len(c.Children) {
		c.Children[c.FocusedIndex].SetFocused(b)
	}
}

func (c *Column) DisplaySize(boundsW, boundsH int) (w, h int) {
	rects := c.GetChildRects(Rect{0, 0, boundsW, boundsH})
	if rects == nil {
		return 0, 0
	}
	height := 0
	width := 0
	for i := range rects {
		height += rects[i].H // combined height
		if rects[i].W > width {
			width = rects[i].W // only the maximum width
		}
	}
	return width, height
}

func (c *Column) Draw(rect Rect, s tcell.Screen) {
	rects := c.GetChildRects(rect)
	for i := range rects {
		c.Children[i].Draw(rects[i], s)
	}
}
