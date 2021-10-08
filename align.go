package dos

import (
	"github.com/gdamore/tcell/v2"
)

type Positioning uint8

const (
	// Inherit is the default layout for all Widgets. The Rect property will be
	// ignored. Calling Align's DisplaySize will return DisplaySize on the
	// Child.
	Inherit Positioning = iota
	// Absolute positioning causes a Widget to be placed at any X, Y coordinate
	// with any arbitrary width and height as specified. This is useful for
	// drop-down menus or other floating widgets. Calling Align's DisplaySize
	// will return zero.
	Absolute
	// Relative positioning is similar to Absolute, but causes a Widget to
	// inherit its parent's position. Calling Align's DisplaySize will return
	// zero.
	Relative
)

type Align struct {
	Child       Widget
	Positioning Positioning
	Rect        Rect // Rect of Child if Positioning is Absolute or Relative.
}

func (a *Align) GetChildRect(currentRect Rect) Rect {
	switch a.Positioning {
	case Absolute:
		return a.Rect
	case Relative:
		return Rect{currentRect.X + a.Rect.X, currentRect.Y + a.Rect.Y, a.Rect.W, a.Rect.H}
	default:
		return currentRect
	}
}

func (a *Align) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if a.Child != nil {
		return a.Child.HandleMouse(a.GetChildRect(currentRect), ev)
	}
	return false
}

func (a *Align) HandleKey(ev *tcell.EventKey) bool {
	if a.Child != nil {
		return a.Child.HandleKey(ev)
	}
	return false
}

func (a *Align) SetFocused(b bool) {
	if a.Child != nil {
		a.Child.SetFocused(b)
	}
}

func (a *Align) DisplaySize(boundsW, boundsH int) (w, h int) {
	if a.Child != nil {
		switch a.Positioning {
		case Absolute:
			fallthrough
		case Relative:
			return 0, 0
		default:
			return a.Child.DisplaySize(boundsW, boundsH)
		}
	}
	return 0, 0
}

func (a *Align) Draw(rect Rect, s tcell.Screen) {
	if a.Child != nil {
		a.Child.Draw(a.GetChildRect(rect), s)
	}
}
