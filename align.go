package dos

import (
	"github.com/gdamore/tcell/v2"
)

type Positioning uint8

const (
	// Boxed is the default layout for all Widgets.
	Boxed Positioning = iota
	// Absolute positioning causes a Widget to be placed at any X, Y coordinate
	// with any arbitrary width and height as specified. This is useful for
	// drop-down menus or other floating widgets.
	Absolute
	// Relative positioning is similar to Absolute, but causes a Widget to inherit its parent's position.
	Relative // TODO: implement relative positioning
)

type Align struct {
	Child       Widget
	Positioning Positioning
	Rect        Rect // Rect of Child if Positioning is Absolute or Relative.
}

func (a *Align) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if a.Child != nil {
		if a.Positioning == Absolute {
			return a.Child.HandleMouse(a.Rect, ev)
		} else {
			return a.Child.HandleMouse(currentRect, ev)
		}
	}
	return false
}

func (a *Align) HandleKey(ev *tcell.EventKey) bool {
	if a.Child != nil {
		if a.Positioning == Absolute {
			return a.Child.HandleKey(ev)
		} else {
			return a.Child.HandleKey(ev)
		}
	}
	return false
}

func (a *Align) SetFocused(b bool) {
	if a.Child != nil {
		a.Child.SetFocused(b)
	}
}

func (a *Align) DisplaySize(boundsW, boundsH int) (w, h int) {
	if a.Positioning != Absolute && a.Child != nil {
		return a.Child.DisplaySize(boundsW, boundsH)
	}
	return 0, 0
}

func (a *Align) Draw(rect Rect, s tcell.Screen) {
	if a.Child != nil {
		if a.Positioning == Absolute {
			a.Child.Draw(a.Rect, s)
		} else {
			a.Child.Draw(rect, s)
		}
	}
}
