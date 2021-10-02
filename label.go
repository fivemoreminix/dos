package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Alignment uint8

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
	AlignJustify
)

type Label struct {
	Text    string
	Align   Alignment
	WrapLen int // Cause the text to wrap after a specified number of terminal cells
	Style   tcell.Style
}

func (l Label) HandleMouse(_ Rect, _ *tcell.EventMouse) bool {
	return false
}

func (l Label) HandleKey(_ *tcell.EventKey) bool {
	return false
}

func (l Label) SetFocused(_ bool) {}

func (l Label) DisplaySize(boundsW, boundsH int) (w, h int) {
	// TODO: account for text wrapping
	return runewidth.StringWidth(l.Text), 1
}

func (l Label) Draw(rect Rect, s tcell.Screen) {
	// TODO: handle overflowing at the edge of rect, overflowing at WrapLen, and text alignment
	DrawString(rect.X, rect.Y, l.Text, l.Style, s)
}
