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

	widget := &dos.Row{
		Children: []dos.Widget{
			&dos.Column{ // Labels
				Children: []dos.Widget{
					&dos.Label{Text: "Username: "},
					&dos.Label{Text: "Password: "},
					&dos.Label{Text: "Favorite number: "},
				},
				HorizontalAlign: dos.AlignRight,
			}, // Fields
			&dos.Column{
				Children: []dos.Widget{
					// TODO: input fields
				},
				OnKeyEvent: func(col *dos.Column, ev *tcell.EventKey) bool {
					if ev.Key() == tcell.KeyTab {
						col.FocusNext()
						return true
					}
					return false
				},
			},
		},
		FocusedIndex: 1,
	}

	var app dos.App
	app = dos.App{
		MainWidget: widget,
		OnKeyEvent: func(ev *tcell.EventKey) bool {
			if ev.Key() == tcell.KeyEsc {
				app.Running = false
				return true
			}
			return false
		},
	}
	app.Run(screen)
}
