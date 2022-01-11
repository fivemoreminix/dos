package main

import (
	"fmt"
	"os"

	"github.com/fivemoreminix/dos"
	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tcell screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize: %v", err)
	}

	var app dos.App
	app = dos.App{
		MainWidget: NewMainWidget(),
		OnKeyEvent: func(ev *tcell.EventKey) bool {
			if ev.Key() == tcell.KeyEsc {
				app.Running = false // Stop the app
				return true         // Report we handled the event
			}
			return false
		},
	}
	app.Run(screen)
}

type MainWidget struct {
	align dos.Align
	label *dos.Label
}

func NewMainWidget() *MainWidget {
	label := &dos.Label{
		Text: "Type something!",
	}
	return &MainWidget{
		align: dos.Align{
			Child: &dos.Box{
				Child: label,
			},
			Positioning: dos.Absolute,
		},
		label: label,
	}
}

func (m *MainWidget) HandleMouse(currentRect dos.Rect, ev *tcell.EventMouse) bool {
	curX, curY := ev.Position()
	m.align.Rect.X, m.align.Rect.Y = curX-m.align.Rect.W, curY-m.align.Rect.H
	return true
}

func (m *MainWidget) HandleKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyBS || ev.Key() == tcell.KeyDEL {
		// Delete the character at the end
		if len(m.label.Text) > 0 {
			m.label.Text = m.label.Text[:len(m.label.Text)-1]
		}
	} else {
		// Insert the typed character at the end
		m.label.Text = string(append([]rune(m.label.Text), ev.Rune()))
	}
	return true
}

func (m *MainWidget) SetFocused(b bool) {
	m.align.SetFocused(b)
}

func (m *MainWidget) DisplaySize(boundsW, boundsH int) (w, h int) {
	return m.align.DisplaySize(boundsW, boundsH)
}

func (m *MainWidget) Draw(rect dos.Rect, s tcell.Screen) {
	// Because the Align is set to Absolute positioning, we give it a position and size through m.align.Rect
	// I cheat and directly access the align's child to know how big it plans to be, because the align will always
	// return 0,0 for an absolute size (as expected)
	m.align.Rect.W, m.align.Rect.H = m.align.Child.DisplaySize(rect.W-m.align.Rect.X, rect.H-m.align.Rect.Y)
	m.align.Draw(rect, s)
}
