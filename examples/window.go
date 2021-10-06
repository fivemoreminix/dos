//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/fivemoreminix/dos"
	"github.com/gdamore/tcell/v2"
)

var (
	windowStyle = tcell.Style{}.Background(tcell.ColorGrey).Foreground(tcell.ColorBlack)
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
		MainWidget: &dos.Scaffold{
			MenuBar:    nil,
			MainWidget: nil,
			Floating: []dos.Widget{MakeDialog(
				"Hello, world!",
				dos.Rect{5, 3, 30, 5},
				&dos.Center{
					Child: &dos.Column{
						Children: []dos.Widget{
							&dos.Label{
								Text:  "Hello, world!",
								Style: windowStyle,
							},
							&dos.Shadow{
								Child: &dos.Button{
									Text:         "Press me",
									NormalStyle:  windowStyle.Background(tcell.ColorWhite),
									FocusedStyle: windowStyle.Background(tcell.ColorWhite),
									OnPressed:    func() {},
								},
								Style:     tcell.Style{}.Background(tcell.ColorBlack).Foreground(tcell.ColorGray),
								MakeSmall: true,
							},
						},
						HorizontalAlign: dos.AlignLeft,
					},
				},
			)},
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

func MakeDialog(title string, rect dos.Rect, child dos.Widget) dos.Widget {
	align := &dos.Align{
		Child:       nil,
		Positioning: dos.Absolute,
		Rect:        rect,
	}
	align.Child = &dos.Window{
		Title:            title,
		Child:            child,
		HideClose:        true,
		OnClosed:         nil,
		DisableMoving:    false,
		OnMove:           func(posX, posY int) { align.Rect.X = posX; align.Rect.Y = posY },
		CloseButtonStyle: tcell.Style{}.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
		TitleBarStyle:    windowStyle.Background(tcell.ColorWhite),
		WindowStyle:      tcell.Style{}.Background(tcell.ColorGrey).Foreground(tcell.ColorBlack),
	}
	return align
}
