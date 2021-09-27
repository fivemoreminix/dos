package dos

import (
	"github.com/gdamore/tcell/v2"
)

type Center struct {
	Child Widget
}

func (c *Center) HandleClick(ev *tcell.EventMouse) bool {
	if c.Child != nil {
		return c.Child.HandleClick(ev)
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

func (c *Center) DisplaySizeInBounds(boundsW, boundsH int) (w, h int) {
	return boundsW, boundsH
}

func (c *Center) Draw(rect Rect, s tcell.Screen) {
	if c.Child != nil {
		childWidth, childHeight := c.Child.DisplaySizeInBounds(rect.W, rect.H)
		x, y := rect.W/2-childWidth/2, rect.H/2-childHeight/2
		c.Child.Draw(Rect{x, y, childWidth, childHeight}, s)
	}
}
