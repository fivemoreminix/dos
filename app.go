package dos

import (
	"github.com/gdamore/tcell/v2"
)

type App struct {
	MainWidget      Widget
	CustomEventLoop func(app *App, s tcell.Screen)
	OnResize        func(width, height int)
}

func (app *App) Run(s tcell.Screen) {
	s.EnableMouse()
	s.EnablePaste()
	app.MainWidget.SetFocused(true)
	if app.CustomEventLoop != nil {
		app.CustomEventLoop(app, s)
	} else {
		DefaultEventLoop(app, s)
	}
}

func DefaultEventLoop(app *App, s tcell.Screen) {
	w, h := s.Size()
	for {
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
			_ = app.MainWidget.HandleKey(ev)
		case *tcell.EventMouse:
			_ = app.MainWidget.HandleMouse(rect, ev)
		default:
			break
		}
	}
}
