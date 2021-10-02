//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/fivemoreminix/dos"
	"github.com/gdamore/tcell/v2"
	"os"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tcell screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize: %v", err)
	}

	// Here we construct the GUI
	widget := MainWidget{
		Child: &dos.Center{
			Child: &dos.Label{
				Text:  "Hello, world!\nTry resizing the screen!",
				Style: tcell.StyleDefault.Foreground(tcell.ColorRed),
			},
		},
	}

	app := dos.App{
		MainWidget: &widget,
	}
	app.Run(screen)
}

// The MainWidget sole purpose is to be the root of the application and capture all input ahead
// of everything else, so we can check if we need to exit.
type MainWidget struct {
	Child dos.Widget
	s     tcell.Screen
}

func (m *MainWidget) HandleMouse(currentRect dos.Rect, ev *tcell.EventMouse) bool {
	if m.Child != nil { // Safely forward input to our child if they exist
		return m.Child.HandleMouse(currentRect, ev)
	}
	return false
}

func (m *MainWidget) HandleKey(ev *tcell2.EventKey) bool {
	if ev.Key() == tcell.KeyEsc {
		m.s.Fini() // Reset the state of the terminal
		os.Exit(0)
	}
	// If we couldn't handle the event, we pass it on to our child
	if m.Child != nil {
		return m.Child.HandleKey(ev)
	}
	return false
}

func (m *MainWidget) SetFocused(b bool) {
	if m.Child != nil {
		m.Child.SetFocused(b)
	}
}

func (m *MainWidget) DisplaySize(boundsW, boundsH int) (w, h int) {
	// If we have a child, we are their size, otherwise we have no size.
	if m.Child != nil {
		return m.Child.DisplaySize(boundsW, boundsH)
	}
	return 0, 0
}

func (m *MainWidget) Draw(rect dos.Rect, s tcell.Screen) {
	// This is a hack to use our screen in other functions (which is really just useful for exiting)
	if m.s == nil {
		m.s = s
	}

	// Now we draw our child in the place of our widget
	if m.Child != nil {
		m.Child.Draw(rect, s)
	}
}
