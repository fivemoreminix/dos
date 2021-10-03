package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Window struct {
	Rect             Rect
	Title            string
	Child            Widget
	DisableClose     bool
	OnClosed         func()
	DisableMoving    bool
	DisableResizing  bool
	CloseButtonStyle tcell.Style
	TitleBarStyle    tcell.Style
	PanelStyle       tcell.Style

	// Window dragging-related variables
	dragging           bool // True if previously received click event on title
	relativeCursorPosX int  // Cursor X position when dragging became true
	relativeCursorPosY int  // Cursor Y position when dragging became true
}

func (w *Window) Close() {
	if w.OnClosed != nil {
		w.OnClosed()
	}
}

func (w *Window) GetChildRect() *Rect {
	if w.Child != nil {
		childW, childH := w.Child.DisplaySize(w.Rect.W, w.Rect.H-1)
		return &Rect{w.Rect.X, w.Rect.Y + 1, childW, childH}
	}
	return nil
}

func (w *Window) HandleMouse(_ Rect, ev *tcell.EventMouse) bool {
	posX, posY := ev.Position()

	if w.dragging {
		if ev.Buttons()&tcell.ButtonPrimary != 0 {
			w.dragging = false // Stop dragging if another click event comes
		}
		w.Rect.X = posX - w.relativeCursorPosX
		w.Rect.Y = posY - w.relativeCursorPosY
		w.SetFocused(true)
		return true
	}

	if w.Child != nil {
		if w.Child.HandleMouse(*w.GetChildRect(), ev) {
			w.SetFocused(true)
		}
	}

	if ev.Buttons()&tcell.ButtonPrimary != 0 {
		// The X button is 3 cells wide, one tall at the top left
		if !w.DisableClose && posX >= w.Rect.X && posX <= w.Rect.X+3 && posY == w.Rect.Y {
			w.SetFocused(true)
			w.Close()
			return true
		}

		// User clicked anywhere on the title bar
		if !w.DisableMoving && posX >= w.Rect.X && posX < w.Rect.X+w.Rect.W && posY == w.Rect.Y {
			w.dragging = true
			w.relativeCursorPosX = posX - w.Rect.X
			w.relativeCursorPosY = 0
			w.SetFocused(true)
			return true
		}

		// User clicked somewhere inside the window
		if w.Rect.HasPoint(posX, posY) {
			w.SetFocused(true) // We're definitely focused after being clicked on
			//if w.Child != nil {
			//	// Maybe the event was meant for our child
			//	childW, childH := w.Child.DisplaySize(w.Rect.W, w.Rect.H)
			//	_ = w.Child.HandleMouse(Rect{w.Rect.X, w.Rect.Y, childW, childH}, ev)
			//}
			return true // Return true because we did "handle" the event regardless
		}
	}
	return false
}

func (w *Window) HandleKey(ev *tcell.EventKey) bool {
	if w.Child != nil {
		return w.Child.HandleKey(ev)
	}
	return false
}

func (w *Window) SetFocused(b bool) {
	if w.Child != nil {
		w.Child.SetFocused(b)
	}
}

func (w *Window) DisplaySize(boundsW, boundsH int) (int, int) {
	return w.Rect.W, w.Rect.H
}

func (w *Window) Draw(_ Rect, s tcell.Screen) {
	for col := 0; col < w.Rect.W; col++ {
		s.SetContent(w.Rect.X+col, w.Rect.Y, ' ', nil, w.TitleBarStyle)
		for row := 1; row < w.Rect.H; row++ {
			s.SetContent(w.Rect.X+col, w.Rect.Y+row, ' ', nil, w.PanelStyle)
		}
	}
	// Draw title
	titleWidth := runewidth.StringWidth(w.Title)
	col := w.Rect.W/2 - titleWidth/2 // Center title
	DrawString(w.Rect.X+col, w.Rect.Y, w.Title, w.TitleBarStyle, s)
	// Draw close button
	DrawString(w.Rect.X, w.Rect.Y, " X ", w.CloseButtonStyle, s)
	// Draw child
	if w.Child != nil {
		w.Child.Draw(*w.GetChildRect(), s)
	}
}
