//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/fivemoreminix/dos"
	"github.com/gdamore/tcell/v2"
	"os"
)

var theme = map[string]tcell.Style{
	"menubar":         tcell.Style{}.Background(tcell.ColorLightGrey).Foreground(tcell.ColorBlack),
	"menubar.focused": tcell.Style{}.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack),
}
var menuBoxDecoration = dos.DefaultBoxDecoration.WithStyle(theme["menubar"])

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tcell screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize: %v", err)
	}
	var app dos.App

	windows := make([]dos.Widget, 0, 3)

	scaffold := &dos.Scaffold{
		MenuBar: &dos.MenuBar{
			Menus: []dos.MenuBarItem{
				makeMenuBarItem("File", []dos.MenuItem{
					dos.MenuItem{
						Title:  "Open file",
						Action: func() {},
					},
					dos.MenuItem{
						Title:  "Save as...",
						Action: func() {},
					},
					dos.MenuItem{Type: dos.MenuItemSeparator},
					dos.MenuItem{
						Title:  "Exit",
						Action: func() { app.Running = false },
					},
				}),
				makeMenuBarItem("Edit", []dos.MenuItem{
					dos.MenuItem{
						Title:  "Copy",
						Action: func() {},
					},
					dos.MenuItem{
						Title:  "Cut",
						Action: func() {},
					},
					dos.MenuItem{Type: dos.MenuItemSeparator},
					dos.MenuItem{
						Title:  "Paste",
						Action: func() { app.Running = false },
					},
				}),
			},
			NormalStyle:    theme["menubar"],
			SelectionStyle: theme["menubar.focused"],
			Selected:       0,
		},
		MainWidget: nil,
		Floating:   windows,
	}

	app = dos.App{
		MainWidget: scaffold,
		OnKeyEvent: func(ev *tcell.EventKey) bool {
			if ev.Key() == tcell.KeyEsc {
				if scaffold.IsMenuBarFocused() {
					if len(scaffold.Floating) > 0 {
						scaffold.FocusFloating()
					} else {
						scaffold.FocusMainWidget()
					}
				} else {
					scaffold.FocusMenuBar()
				}
				return true
			}
			return false
		},
	}
	app.Run(screen)
}

func makeMenuBarItem(title string, items []dos.MenuItem) dos.MenuBarItem {
	return dos.MenuBarItem{
		Title: title,
		Menu: dos.Menu{
			Items:          items,
			Decorated:      true,
			Decoration:     &menuBoxDecoration,
			NormalStyle:    theme["menubar"],
			SelectionStyle: theme["menubar.focused"],
		},
	}
}
