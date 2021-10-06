package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// A Window displays a flexible dialog. The dialog can be moved by the user,
// but events must be handled with the OnMove and OnClosed callbacks. By
// standard, the Window assumes the position and size of whatever Rect is
// provided to Draw and DisplaySize, but you could use an Align to provide the
// Window an arbitrary position and size on the terminal.
type Window struct {
	Title            string
	Child            Widget
	HideClose        bool
	OnClosed         func()
	DisableMoving    bool
	OnMove           func(newX, newY int)
	CloseButtonStyle tcell.Style
	TitleBarStyle    tcell.Style
	WindowStyle      tcell.Style

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

func (w *Window) GetChildRect(rect Rect) *Rect {
	if w.Child != nil {
		childW, childH := w.Child.DisplaySize(rect.W, rect.H-1)
		return &Rect{rect.X, rect.Y + 1, childW, childH}
	}
	return nil
}

func (w *Window) HandleMouse(rect Rect, ev *tcell.EventMouse) bool {
	posX, posY := ev.Position()

	if w.dragging {
		if ev.Buttons()&tcell.ButtonPrimary != 0 {
			w.dragging = false // Stop dragging if another click event comes
		}
		w.SetFocused(true)
		// TODO: pass a relative movement difference
		if w.OnMove != nil {
			w.OnMove(posX-w.relativeCursorPosX, posY-w.relativeCursorPosY)
		}
		return true
	}

	if w.Child != nil {
		if w.Child.HandleMouse(*w.GetChildRect(rect), ev) {
			w.SetFocused(true)
		}
	}

	if ev.Buttons()&tcell.ButtonPrimary != 0 {
		// The X button is 3 cells wide, one tall at the top left
		if !w.HideClose && posX >= rect.X && posX <= rect.X+3 && posY == rect.Y {
			w.SetFocused(true)
			w.Close()
			return true
		}

		// User clicked anywhere on the title bar
		if !w.DisableMoving && posX >= rect.X && posX < rect.X+rect.W && posY == rect.Y {
			w.dragging = true
			w.relativeCursorPosX = posX - rect.X
			w.relativeCursorPosY = 0
			w.SetFocused(true)
			return true
		}

		// User clicked somewhere inside the window
		if rect.HasPoint(posX, posY) {
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
	return boundsW, boundsH
}

func (w *Window) Draw(rect Rect, s tcell.Screen) {
	for col := 0; col < rect.W; col++ {
		s.SetContent(rect.X+col, rect.Y, ' ', nil, w.TitleBarStyle)
		for row := 1; row < rect.H; row++ {
			s.SetContent(rect.X+col, rect.Y+row, ' ', nil, w.WindowStyle)
		}
	}
	// Draw title
	titleWidth := runewidth.StringWidth(w.Title)
	col := rect.W/2 - titleWidth/2 // Center title
	DrawString(rect.X+col, rect.Y, w.Title, w.TitleBarStyle, s)
	// Draw close button
	if !w.HideClose {
		DrawString(rect.X, rect.Y, " X ", w.CloseButtonStyle, s)
	}
	// Draw child
	if w.Child != nil {
		w.Child.Draw(*w.GetChildRect(rect), s)
	}
}
