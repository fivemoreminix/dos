package dos

import "github.com/gdamore/tcell/v2"

type TextInput struct {
	Text             string // User-entered content.
	Placeholder      string // Placeholder is visible when Text is empty.
	IsHidden         bool
	HiddenChar       rune
	Scroll           int // Number of runes skipped when viewing.
	Width            int // If Width is zero, then it is will be as wide as possible.
	NormalStyle      tcell.Style
	FocusedStyle     tcell.Style
	PlaceholderStyle tcell.Style // If PlaceholderStyle is zero (or default style), then it inherits Normal/FocusedStyle.
	OnTextEdited     func(text string)

	cursorPos int
	focused   bool
}

func (t *TextInput) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	panic("implement me")
}

func (t *TextInput) HandleKey(ev *tcell.EventKey) bool {
	panic("implement me")
}

func (t *TextInput) SetFocused(b bool) {
	t.focused = b
}

func (t *TextInput) DisplaySize(boundsW, boundsH int) (w, h int) {
	panic("implement me")
}

func (t *TextInput) Draw(rect Rect, s tcell.Screen) {
	panic("implement me")
}
