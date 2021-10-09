package dos

import (
	"github.com/gdamore/tcell/v2"
)

// The BoxDecoration allows for an individual rune per side and corner of the
// box being drawn. See https://en.wikipedia.org/wiki/Box-drawing_character
type BoxDecoration struct {
	Hor, Vert      rune // Horizontal and vertical sides
	TL, TR, BR, BL rune // Clockwise corners
	JointT, JointR rune // Joints (for menus and other divided boxes)
	JointB, JointL rune
	Style          tcell.Style
}

// WithStyle is a helper function to make a copy of the BoxDecoration with a
// different style.
func (b BoxDecoration) WithStyle(style tcell.Style) BoxDecoration {
	b.Style = style
	return b
}

// A Box draws an enclosed rectangle around its Child. A Box assumes each rune
// has a width of one terminal cell so double-wide characters will not be drawn
// correctly using a Box and BoxDecoration combination.
type Box struct {
	Child      Widget
	Decoration *BoxDecoration
}

var DefaultBoxDecoration = BoxDecoration{
	Hor:    '─',
	Vert:   '│',
	TL:     '┌',
	TR:     '┐',
	BR:     '┘',
	BL:     '└',
	JointT: '┬',
	JointR: '┤',
	JointB: '┴',
	JointL: '├',
	Style:  tcell.StyleDefault,
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

	DrawBox(rect, decoration, s)

	if b.Child != nil {
		b.Child.Draw(Rect{rect.X + 1, rect.Y + 1, rect.W - 1, rect.H - 1}, s)
	}
}
