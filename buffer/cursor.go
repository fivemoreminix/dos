package buffer

import (
	"math"
	"unicode"
)

// So why is the code for moving the cursor in the buffer package, and not in the
// TextEdit component? Well, it used to be, but it sucked that way. The cursor
// needs to have a reference to the buffer to know where lines end and how it can
// move. The Buffer is the city, and the Cursor is the car.

// A Selection represents a region of the buffer to be selected for text editing
// purposes. It is asserted that the start position is less than the end position.
// The start and end are inclusive. If the EndCol of a Region is one more than the
// last column of a line, then it points to the line delimiter at the end of that
// line. It is understood that as a Region spans multiple lines, those connecting
// line-delimiters are included in the selection, as well.
type Region struct {
	Start *Cursor
	End   *Cursor
}

func NewRegion(in Buffer) Region {
	return Region{
		NewCursor(in),
		NewCursor(in),
	}
}

// A Cursor's functions emulate common cursor actions. To have a Cursor be
// automatically updated when the buffer has text prepended or appended -- one
// should register the Cursor with the Buffer's function `RegisterCursor()`
// which makes the Cursor "anchored" to the Buffers contents when they change.
type Cursor struct {
	buffer  Buffer
	prevCol int
	Line    int
	Col     int
}

func NewCursor(in Buffer) *Cursor {
	return &Cursor{
		buffer: in,
	}
}

func (c *Cursor) Left() {
	if c.Col == 0 && c.Line != 0 { // If we are at the beginning of the current line...
		// Go to the end of the above line
		c.Line--
		c.Col, _ = c.buffer.RunesInLine(c.Line, false)
	} else {
		c.Col = Max(c.Col-1, 0)
	}
	c.prevCol = c.Col
}

func (c *Cursor) Right() {
	runes, _ := c.buffer.RunesInLine(c.Line, false)
	if c.Col >= runes && c.Line < c.buffer.Lines()-1 {
		// If we are at the end of the current line,
		// and not at the last line...
		c.Line, c.Col = c.buffer.ClampLineCol(c.Line+1, 0) // Go to beginning of line below
	} else {
		c.Line, c.Col = c.buffer.ClampLineCol(c.Line, c.Col+1)
	}
	c.prevCol = c.Col
}

func (c *Cursor) Up() {
	if c.Line == 0 { // If the cursor is at the first line...
		c.Line, c.Col = 0, 0 // Go to beginning
	} else {
		c.Line, c.Col = c.buffer.ClampLineCol(c.Line-1, c.Col)
	}
}

func (c *Cursor) Down() {
	if c.Line == c.buffer.Lines()-1 { // If the cursor is at the last line...
		c.Line, c.Col = c.buffer.ClampLineCol(c.Line, math.MaxInt32) // Go to end of current line
	} else {
		c.Line, c.Col = c.buffer.ClampLineCol(c.Line+1, c.Col)
	}
}

// NextWordBoundaryEnd proceeds to the position after the last character of the
// next word boundary to the right of the Cursor. A word boundary is the
// beginning or end of any sequence of similar or same-classed characters.
// Whitespace is skipped.
func (c *Cursor) NextWordBoundaryEnd() {
	// Get position of cursor in buffer as pos
	// get classification of character at pos or assume none if whitespace
	// for each pos until end of buffer: pos + 1 (at end)
	//	 if pos char is not of previous pos char class:
	//      set cursor position as pos
	//

	// only skip contiguous characters for word characters
	// jump to position *after* any symbols

	pos := c.buffer.LineColToPos(c.Line, c.Col)
	r, _ := c.buffer.RuneAtPos(pos)
	startClass := getRuneCharclass(r)
	pos++
	c.buffer.EachRuneFromPos(pos, func(rpos int, r rune) bool {
		class := getRuneCharclass(r)
		if class != startClass && class != charwhitespace {
			return true
		}
		return false
	})

	c.Line, c.Col = c.buffer.PosToLineCol(pos)
}

func (c *Cursor) PrevWordBoundaryStart() {
	// TODO
}

// LineCol sets the Line and Col of the Cursor to those provided. `line` is
// clamped within the range [0, lines in buffer). `col` is then clamped within
// the range [0, line length in runes minus delimiter).
func (c *Cursor) LineCol(line, col int) {
	c.Line, c.Col = c.buffer.ClampLineCol(line, col)
}

func (c *Cursor) Eq(other *Cursor) bool {
	return c.buffer == other.buffer && c.Line == other.Line && c.Col == other.Col
}

type charclass uint8

const (
	charwhitespace charclass = iota
	charword
	charsymbol
)

func getRuneCharclass(r rune) charclass {
	if unicode.IsSpace(r) {
		return charwhitespace
	} else if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
		return charword
	} else {
		return charsymbol
	}
}
