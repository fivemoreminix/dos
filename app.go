package dos

import (
	"github.com/gdamore/tcell/v2"
)

type App struct {
	MainWidget      Widget
	CustomEventLoop func(app *App, s tcell.Screen)
	Running         bool
	OnResize        func(width, height int)
	// OnKeyEvent is called before the MainWidget's handler, and if this function
	// returns true, then the event is never passed onto the main widget.
	OnKeyEvent func(ev *tcell.EventKey) bool
	// OnMouseEvent is called before the MainWidget's handler, and if this function
	// returns true, then the event is never passed onto the main widget.
	OnMouseEvent func(ev *tcell.EventMouse) bool
}

func (app *App) Run(s tcell.Screen) {
	s.EnableMouse()
	s.EnablePaste()
	app.MainWidget.SetFocused(true)
	app.Running = true
	if app.CustomEventLoop != nil {
		app.CustomEventLoop(app, s)
	} else {
		DefaultEventLoop(app, s)
	}
	s.Fini()
}

func DefaultEventLoop(app *App, s tcell.Screen) {
	w, h := s.Size()
	for app.Running {
		s.Clear()

		rect := Rect{0, 0, w, h}
		app.MainWidget.Draw(rect, s)

		s.Show() // Renders all changed cells

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			w, h = ev.Size()
			if app.OnResize != nil {
				app.OnResize(w, h)
			}
			s.Sync() // Redraw the entire screen
		case *tcell.EventKey:
			if app.OnKeyEvent != nil {
				if app.OnKeyEvent(ev) {
					break
				}
			}
			_ = app.MainWidget.HandleKey(ev)
		case *tcell.EventMouse:
			if app.OnMouseEvent != nil {
				if app.OnMouseEvent(ev) {
					break
				}
			}
			_ = app.MainWidget.HandleMouse(rect, ev)
		}
	}
}
