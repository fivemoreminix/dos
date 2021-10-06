package dos

import (
	"github.com/gdamore/tcell/v2"
)

type Center struct {
	Child Widget
}

func (c *Center) GetChildRect(currentRect Rect) Rect {
	childWidth, childHeight := c.Child.DisplaySize(currentRect.W, currentRect.H)
	// Here I subtract 1 from both the x and the y because the widths & heights
	// are 1-based, while positions are zero-based. Without adding one, all
	// centered positions are right one, and down one of their correct locations.
	x, y := currentRect.W/2-1-childWidth/2, currentRect.H/2-1-childHeight/2
	return Rect{currentRect.X + x, currentRect.Y + y, childWidth, childHeight}
}

func (c *Center) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if c.Child != nil {
		return c.Child.HandleMouse(c.GetChildRect(currentRect), ev)
	} else {
		return false
	}
}

func (c *Center) HandleKey(ev *tcell.EventKey) bool {
	if c.Child != nil {
		return c.Child.HandleKey(ev)
	} else {
		return false
	}
}

func (c *Center) SetFocused(b bool) {
	if c.Child != nil {
		c.Child.SetFocused(b)
	}
}

func (c *Center) DisplaySize(boundsW, boundsH int) (w, h int) {
	return boundsW, boundsH
}

func (c *Center) Draw(rect Rect, s tcell.Screen) {
	if c.Child != nil {
		c.Child.Draw(c.GetChildRect(rect), s)
	}
}
