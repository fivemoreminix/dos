package dos

import (
	"github.com/gdamore/tcell/v2"
)

// The BoxDecoration allows for an individual rune per side and corner of the
// box being drawn. See https://en.wikipedia.org/wiki/Box-drawing_character
type BoxDecoration struct {
	T, R, B, L     rune // Clockwise sides
	TL, TR, BR, BL rune // Clockwise corners
	Style          tcell.Style
}

// A Box draws an enclosed rectangle around its Child. A Box assumes each rune
// has a width of one terminal cell so double-wide characters will not be drawn
// correctly using a Box and BoxDecoration combination.
type Box struct {
	Child      Widget
	Decoration *BoxDecoration
}

var DefaultBoxDecoration = BoxDecoration{
	T:     '─',
	R:     '│',
	B:     '─',
	L:     '│',
	TL:    '┌',
	TR:    '┐',
	BR:    '┘',
	BL:    '└',
	Style: tcell.StyleDefault,
}

func (b *Box) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if b.Child != nil {
		return b.Child.HandleMouse(currentRect, ev)
	}
	return false
}

func (b *Box) HandleKey(ev *tcell.EventKey) bool {
	if b.Child != nil {
		return b.Child.HandleKey(ev)
	}
	return false
}

func (b *Box) SetFocused(v bool) {
	if b.Child != nil {
		b.Child.SetFocused(v)
	}
}

func (b *Box) DisplaySize(boundsW, boundsH int) (w, h int) {
	if b.Child != nil {
		childW, childH := b.Child.DisplaySize(boundsW-2, boundsH-2)
		return childW + 2, childH + 2
	}
	return 0, 0
}

func (b *Box) Draw(rect Rect, s tcell.Screen) {
	// Do not draw if not even a single cell of room
	if rect.W < 1 || rect.H < 1 {
		return
	}

	decoration := b.Decoration
	if decoration == nil {
		decoration = &DefaultBoxDecoration
	}

	// Draw top left and bottom left corners
	s.SetContent(rect.X, rect.Y, decoration.TL, nil, decoration.Style)
	s.SetContent(rect.X, rect.Y+rect.H-1, decoration.BL, nil, decoration.Style)
	// Draw top right and bottom right corners
	s.SetContent(rect.X+rect.W-1, rect.Y, decoration.TR, nil, decoration.Style)
	s.SetContent(rect.X+rect.W-1, rect.Y+rect.H-1, decoration.BR, nil, decoration.Style)
	// Draw top and bottom sides
	for col := 1; col < rect.W-1; col++ {
		s.SetContent(rect.X+col, rect.Y, decoration.T, nil, decoration.Style)
		s.SetContent(rect.X+col, rect.Y+rect.H-1, decoration.B, nil, decoration.Style)
	}
	// Draw left and right sides
	for row := 1; row < rect.H-1; row++ {
		s.SetContent(rect.X, rect.Y+row, decoration.L, nil, decoration.Style)
		s.SetContent(rect.X+rect.W-1, rect.Y+row, decoration.R, nil, decoration.Style)
	}

	if b.Child != nil {
		b.Child.Draw(Rect{rect.X + 1, rect.Y + 1, rect.W - 1, rect.H - 1}, s)
	}
}
