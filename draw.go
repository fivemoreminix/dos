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

func DrawBox(rect Rect, decoration *BoxDecoration, s tcell.Screen) {
	// Draw top left and bottom left corners
	s.SetContent(rect.X, rect.Y, decoration.TL, nil, decoration.Style)
	s.SetContent(rect.X, rect.Y+rect.H-1, decoration.BL, nil, decoration.Style)
	// Draw top right and bottom right corners
	s.SetContent(rect.X+rect.W-1, rect.Y, decoration.TR, nil, decoration.Style)
	s.SetContent(rect.X+rect.W-1, rect.Y+rect.H-1, decoration.BR, nil, decoration.Style)
	// Draw top and bottom sides
	for col := 1; col < rect.W-1; col++ {
		s.SetContent(rect.X+col, rect.Y, decoration.Hor, nil, decoration.Style)
		s.SetContent(rect.X+col, rect.Y+rect.H-1, decoration.Hor, nil, decoration.Style)
	}
	// Draw left and right sides
	for row := 1; row < rect.H-1; row++ {
		s.SetContent(rect.X, rect.Y+row, decoration.Vert, nil, decoration.Style)
		s.SetContent(rect.X+rect.W-1, rect.Y+row, decoration.Vert, nil, decoration.Style)
	}
}
