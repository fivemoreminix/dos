//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/fivemoreminix/dos"
	"github.com/gdamore/tcell/v2"
	"os"
)

const (
	areYouReady = "你准备好了吗?"
)

var styleDefault = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tcell screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize: %v", err)
	}
	defer screen.Fini()

	width, height := screen.Size()

	var app dos.App

	windowStyle := styleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite)
	windowTitleBarStyle := windowStyle.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	windows := make([]dos.Window, 0, 6)
	newWindow := func(window dos.Window) {
		// lazy architecture for the demo
		windows := &app.MainWidget.(*dos.Scaffold).Windows
		if len(*windows) > 0 {
			(*windows)[len(*windows)-1].SetFocused(false)
		}
		*windows = append(*windows, window)
		(*windows)[len(*windows)-1].SetFocused(true)
	}
	closeFocusedWindow := func() {
		windows := &app.MainWidget.(*dos.Scaffold).Windows
		*windows = (*windows)[:len(*windows)-1]
		if len(*windows) > 0 {
			(*windows)[len(*windows)-1].SetFocused(true)
		}
	}

	defaultWindow := dos.Window{
		Rect:  dos.Rect{3, 3, 30, 20},
		Title: "窗口对话",
		Child: &dos.Button{
			Text:         "Click me...",
			NormalStyle:  windowStyle.Background(tcell.ColorGrey),
			FocusedStyle: windowStyle.Background(tcell.ColorGrey),
			OnPressed: func() {
				newWindow(dos.Window{
					Rect:  dos.Rect{width/2 - 30/2, height / 3, 30, 5},
					Title: "Dialog",
					Child: &dos.Label{
						Text:  "You clicked the button.",
						Style: windowStyle.Foreground(tcell.ColorRed),
					},
					OnClosed:         closeFocusedWindow,
					CloseButtonStyle: windowStyle.Background(tcell.ColorRed),
					TitleBarStyle:    windowTitleBarStyle,
					PanelStyle:       windowStyle,
				})
			},
		},
		OnClosed:         closeFocusedWindow,
		CloseButtonStyle: windowStyle.Background(tcell.ColorRed),
		TitleBarStyle:    windowTitleBarStyle,
		PanelStyle:       windowStyle,
	}
	windows = append(windows, defaultWindow)

	app = dos.App{
		MainWidget: &dos.Scaffold{
			MenuBar: &dos.MenuBar{
				Menus: []dos.MenuBarItem{
					dos.MenuBarItem{
						Title: "File",
						Menu: dos.Menu{
							Items: []dos.MenuItem{
								dos.MenuItem{
									Title:  "New",
									Type:   dos.MenuItemAction,
									Action: func() { newWindow(defaultWindow) },
								},
								dos.MenuItem{Type: dos.MenuItemSeparator},
								dos.MenuItem{
									Title:  "Exit",
									Type:   dos.MenuItemAction,
									Action: func() { app.Running = false },
								},
							},
							NormalStyle:    windowStyle,
							SelectionStyle: windowStyle.Background(tcell.ColorGrey),
						},
					},
					dos.MenuBarItem{
						Title: "Edit",
						Menu:  dos.Menu{},
					},
					dos.MenuBarItem{
						Title: "Sudoku",
						Menu:  dos.Menu{},
					},
					dos.MenuBarItem{
						Title: "Help?",
						Menu:  dos.Menu{},
					},
				},
				NormalStyle:    styleDefault.Background(tcell.ColorMaroon).Foreground(tcell.ColorWhite),
				SelectionStyle: styleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
			},
			MainWidget: &dos.Center{Child: dos.Label{Text: areYouReady}},
			Windows:    windows,
		},
		OnKeyEvent: func(ev *tcell.EventKey) bool {
			if ev.Key() == tcell.KeyEsc {
				app.Running = false
				return true
			}
			return false
		},
		OnResize: func(w, h int) {
			width, height = w, h
		},
	}
	app.Run(screen)
}
