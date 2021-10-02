package dos

import (
	"github.com/gdamore/tcell/v2"
)

// A Column orders its children vertically.
type Column struct {
	Children        []Widget
	HorizontalAlign Alignment
}

func (c Column) GetChildRects(currentRect Rect) []Rect {
	if c.Children != nil {
		rects := make([]Rect, 0, len(c.Children))
		var rectsHeightSum int
		for i := range c.Children {
			childW, childH := c.Children[i].DisplaySize(currentRect.W, currentRect.H/len(c.Children))
			// Allow child rects to overflow our currentRect, but later we resize them to fit
			rects = append(rects, Rect{
				currentRect.X + currentRect.W/2 - childW/2, // center child in X axis
				rectsHeightSum,
				childW,
				childH,
			})
			rectsHeightSum += childH
		}
		// Check if we overflowed the currentRect height
		if rectsHeightSum > currentRect.H {
			// Gives us the ratio to multiply by
			scaleRatio := float64(currentRect.H) / float64(rectsHeightSum)
			rectsHeightSum = 0
			for i := 0; i < len(rects); i++ {
				r := &rects[i]
				r.Y = rectsHeightSum // Move rect up
				r.H = int(float64(r.H) * scaleRatio)
				rectsHeightSum += r.H
			}
		}
		return rects
	}
	return nil
}

func (c Column) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if c.Children != nil {

	}
	return false
}

func (c Column) HandleKey(ev *tcell.EventKey) bool {
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

func (c Column) DisplaySize(boundsW, boundsH int) (w, h int) {
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
