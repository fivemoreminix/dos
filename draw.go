package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"unicode/utf8"
)

// DrawString prints the string s at column x and row y with the provided style.
// Use mattn/go-runewidth to determine how many terminal cells your string will
// consume.
func DrawString(x, y int, s string, style tcell.Style, screen tcell.Screen) {
	var col int
	var byteIdx int
	for byteIdx < len(s) {
		r, size := utf8.DecodeRuneInString(s[byteIdx:])
		screen.SetContent(x+col, y, r, nil, style)
		byteIdx += size
		col += runewidth.RuneWidth(r)
	}
}

func DrawRect(rect Rect, r rune, style tcell.Style, screen tcell.Screen) {
	for x := rect.X; x < rect.X+rect.W; x++ {
		for y := rect.Y; y < rect.Y+rect.H; y++ {
			screen.SetContent(x, y, r, nil, style)
		}
	}
}
