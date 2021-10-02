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

	windowStyle := styleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite)
	var app dos.App
	app = dos.App{
		MainWidget: &dos.Window{
			Rect:  dos.Rect{3, 3, 30, 20},
			Title: "窗口对话",
			Child: dos.Label{
				Text:  "我是个标签",
				Style: windowStyle,
			},
			DisableClose:     false,
			OnClosed:         nil,
			DisableMoving:    false,
			DisableResizing:  false,
			CloseButtonStyle: windowStyle.Background(tcell.ColorRed),
			TitleBarStyle:    windowStyle.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite),
			PanelStyle:       windowStyle,
		},
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
