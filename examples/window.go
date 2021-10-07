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

	label := &dos.Label{
		Text:    quote,
		Align:   dos.AlignLeft,
		WrapLen: 0,
		Style:   defaultStyle,
	}

	var app dos.App
	app = dos.App{
		ClearStyle: defaultStyle,
		MainWidget: &dos.Scaffold{
			MenuBar:    nil,
			MainWidget: &dos.Center{Child: label},
			Floating:   []dos.Widget{NewMainWindow(label)},
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

type MainWindow struct {
	contentLabel *dos.Label
	child        dos.Widget
}

func NewMainWindow(label *dos.Label) *MainWindow {
	row := &dos.Row{
		Children: []dos.Widget{
			&dos.Padding{
				Child: &dos.Shadow{
					Child: &dos.Button{
						Text:         "Left align",
						NormalStyle:  windowStyle.Background(tcell.ColorWhite),
						FocusedStyle: windowStyle.Background(tcell.ColorWhite),
						OnPressed:    func() { label.Align = dos.AlignLeft },
					},
					Style:     tcell.Style{}.Background(tcell.ColorGray).Foreground(tcell.ColorLightBlue),
					MakeSmall: true,
				},
				Top:    1,
				Right:  1,
				Bottom: 0,
				Left:   1,
			},
			&dos.Padding{
				Child: &dos.Shadow{
					Child: &dos.Button{
						Text:         "Center align",
						NormalStyle:  windowStyle.Background(tcell.ColorWhite),
						FocusedStyle: windowStyle.Background(tcell.ColorWhite),
						OnPressed:    func() { label.Align = dos.AlignCenter },
					},
					Style:     tcell.Style{}.Background(tcell.ColorGray).Foreground(tcell.ColorLightBlue),
					MakeSmall: true,
				},
				Top:    1,
				Right:  1,
				Bottom: 0,
				Left:   1,
			},
			&dos.Padding{
				Child: &dos.Shadow{
					Child: &dos.Button{
						Text:         "Right align",
						NormalStyle:  windowStyle.Background(tcell.ColorWhite),
						FocusedStyle: windowStyle.Background(tcell.ColorWhite),
						OnPressed:    func() { label.Align = dos.AlignRight },
					},
					Style:     tcell.Style{}.Background(tcell.ColorGray).Foreground(tcell.ColorLightBlue),
					MakeSmall: true,
				},
				Top:    1,
				Right:  1,
				Bottom: 0,
				Left:   1,
			},
		},
		VerticalAlign: dos.AlignLeft,
		FocusedIndex:  0,
		OnKeyEvent:    nil,
	}
	row.OnKeyEvent = func(ev *tcell.EventKey) bool {
		if ev.Key() == tcell.KeyTab {
			row.SetFocused(false)
			row.FocusedIndex++
			if row.FocusedIndex >= len(row.Children) {
				row.FocusedIndex = 0
			}
			row.SetFocused(true)
			return true
		}
		return false
	}
	dialog := MakeDialog(
		"Text Alignment",
		dos.Rect{5, 3, 53, 6},
		&dos.Center{
			Child: &dos.Column{
				Children: []dos.Widget{
					&dos.Padding{
						Child: &dos.Label{
							Text:  "Choose a text alignment",
							Style: windowStyle,
						},
						Top: 1,
					},
					row,
				},
				HorizontalAlign: dos.AlignLeft,
				FocusedIndex:    1,
			},
		},
	)
	return &MainWindow{
		contentLabel: label,
		child:        dialog,
	}
}

func (m *MainWindow) HandleMouse(currentRect dos.Rect, ev *tcell.EventMouse) bool {
	return m.child.HandleMouse(currentRect, ev)
}

func (m *MainWindow) HandleKey(ev *tcell.EventKey) bool {
	return m.child.HandleKey(ev)
}

func (m *MainWindow) SetFocused(b bool) {
	m.child.SetFocused(b)
}

func (m *MainWindow) DisplaySize(boundsW, boundsH int) (w, h int) {
	return m.child.DisplaySize(boundsW, boundsH)
}

func (m *MainWindow) Draw(rect dos.Rect, s tcell.Screen) {
	m.child.Draw(rect, s)
}
