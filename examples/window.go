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
	defaultStyle = tcell.Style{}.Background(tcell.ColorBlue).Foreground(tcell.ColorGrey)
	windowStyle  = tcell.Style{}.Background(tcell.ColorLightBlue).Foreground(tcell.ColorBlack)
	quote        = `Are you quite sure that all those bells and whistles,
all those wonderful facilities of your so called powerful programming languages,
belong to the solution set rather than the problem set?

 — Edsger W. Dijkstra`
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
		ClearStyle: defaultStyle,
		MainWidget: &dos.Scaffold{
			MenuBar: nil,
			MainWidget: &dos.Label{
				Text:    quote,
				Align:   dos.AlignLeft,
				WrapLen: 0,
				Style:   defaultStyle,
			},
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
								Style:     tcell.Style{}.Background(tcell.ColorGray).Foreground(tcell.ColorLightBlue),
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
	align.Child = &dos.Shadow{
		Child: &dos.Window{
			Title:            title,
			Child:            child,
			HideClose:        true,
			OnClosed:         nil,
			DisableMoving:    false,
			OnMove:           func(posX, posY int) { align.Rect.X = posX; align.Rect.Y = posY },
			CloseButtonStyle: tcell.Style{}.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
			TitleBarStyle:    windowStyle.Background(tcell.ColorWhite),
			WindowStyle:      windowStyle,
		},
		Style:     defaultStyle.Reverse(true),
		MakeSmall: false,
	}
	return align
}
