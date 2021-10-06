package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"unicode/utf8"
)

type Button struct {
	Text         string
	NormalStyle  tcell.Style
	FocusedStyle tcell.Style
	OnPressed    func()
	focused      bool
}

func (b *Button) Press() {
	if b.OnPressed != nil {
		b.OnPressed()
	}
}

func (b *Button) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if ev.Buttons()&tcell.ButtonPrimary != 0 {
		if currentRect.HasPoint(ev.Position()) {
			b.Press()
			return true
		}
	}
	return false
}

func (b *Button) HandleKey(ev *tcell.EventKey) bool {
	if b.focused && ev.Key() == tcell.KeyEnter || ev.Rune() == ' ' {
		b.Press()
		return true
	}
	return false
}

func (b *Button) SetFocused(v bool) {
	b.focused = v
}

func (b *Button) DisplaySize(boundsW, boundsH int) (w, h int) {
	return runewidth.StringWidth(b.Text) + 4, 1
}

func (b *Button) Draw(rect Rect, s tcell.Screen) {
	w, _ := b.DisplaySize(rect.W, rect.H)

	var style tcell.Style
	if b.focused {
		style = b.FocusedStyle
	} else {
		style = b.NormalStyle
	}

	var col int
	var byteIdx int
	for col < w {
		if col > 1 && col < w-2 { // Draw text
			r, size := utf8.DecodeRune([]byte(b.Text[byteIdx:]))
			s.SetContent(rect.X+col, rect.Y, r, nil, style)
			col += runewidth.RuneWidth(r)
			byteIdx += size
		} else { // Draw padding
			r := ' '
			if b.focused {
				if col == 0 {
					r = '▸'
				} else if col == w-1 {
					r = '◂'
				}
			}
			s.SetContent(rect.X+col, rect.Y, r, nil, style)
			col++
		}
	}
}
