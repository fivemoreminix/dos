package dos

import (
	"github.com/gdamore/tcell/v2"
)

type Center struct {
	Child Widget
}

func (c *Center) GetChildRect(currentRect Rect) Rect {
	childWidth, childHeight := c.Child.DisplaySizeInBounds(currentRect.W, currentRect.H)
	x, y := currentRect.W/2-childWidth/2, currentRect.H/2-childHeight/2
	return Rect{x, y, childWidth, childHeight}
}

func (c *Center) HandleClick(currentRect Rect, ev *tcell.EventMouse) bool {
	if c.Child != nil {
		return c.Child.HandleClick(c.GetChildRect(currentRect), ev)
	} else {
		return false
	}
}

func (c *Center) HandleKey(currentRect Rect, ev *tcell.EventKey) bool {
	if c.Child != nil {
		return c.Child.HandleKey(c.GetChildRect(currentRect), ev)
	} else {
		return false
	}
}

func (c *Center) SetFocused(b bool) {
	if c.Child != nil {
		c.Child.SetFocused(b)
	}
}

func (c *Center) DisplaySizeInBounds(boundsW, boundsH int) (w, h int) {
	return boundsW, boundsH
}

func (c *Center) Draw(rect Rect, s tcell.Screen) {
	if c.Child != nil {
		c.Child.Draw(c.GetChildRect(rect), s)
	}
}
