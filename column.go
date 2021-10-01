package dos

import "github.com/gdamore/tcell/v2"

// A Column orders its children vertically.
type Column struct {
	Children        []Widget
	HorizontalAlign Align
}

func (c Column) GetChildRects(currentRect Rect) []Rect {
	if c.Children != nil {
		rects := make([]Rect, len(c.Children))
		for i := range c.Children {
			childW, childH := c.Children[i].DisplaySizeInBounds(currentRect.W, currentRect.H / len(c.Children))
			// check all child preferred max heights

		}
	}
	return nil
}

func (c Column) HandleClick(currentRect Rect, ev *tcell.EventMouse) bool {
	if c.Children != nil {

	}
	return false
}

func (c Column) HandleKey(currentRect Rect, ev *tcell.EventKey) bool {
	if c.Children != nil {

	}
	return false
}

func (c Column) SetFocused(b bool) {
	if c.Children != nil {
		for i := range c.Children {
			c.Children[i].SetFocused(b)
		}
	}
}

func (c Column) DisplaySizeInBounds(boundsW, boundsH int) (w, h int) {
	rects := c.GetChildRects(Rect{0, 0, boundsW, boundsH})
	if rects == nil {
		return 0, 0
	}
	height := 0
	width := 0
	for i := range rects {
		height += rects[i].H
		if rects[i].W > width {
			width = rects[i].W
		}
	}
	return width, height
	// only as wide and tall as the combined width and height of all its children
}

func (c Column) Draw(rect Rect, s tcell.Screen) {
	panic("implement me")
}
