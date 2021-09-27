package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Label struct {
	Text  string
	Style tcell.Style
}

func (l Label) HandleClick(_ *tcell.EventMouse) bool {
	return false
}

func (l Label) HandleKey(_ *tcell.EventKey) bool {
	return false
}

func (l Label) SetFocused(b bool) {}

func (l Label) DisplaySizeInBounds(boundsW, boundsH int) (w, h int) {
	// TODO: account for text wrapping
	return runewidth.StringWidth(l.Text), 1
}

func (l Label) Draw(rect Rect, s tcell.Screen) {
	if len(l.Text) == 1 {
		// TODO: it is a bug if the label overflows the rect
		s.SetContent(rect.X, rect.Y, rune(l.Text[0]), nil, l.Style)
	} else if len(l.Text) > 1 {
		s.SetContent(rect.X, rect.Y, rune(l.Text[0]), []rune(l.Text[1:]), l.Style)
	}
}
