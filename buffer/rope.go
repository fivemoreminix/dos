package buffer

import (
	bytealg "bytes"
	"io"
	"unicode/utf8"

	ropes "github.com/zyedidia/rope"
)

type RopeBuffer struct {
	rope      *ropes.Node
	anchors   []*Cursor
	lineDelim string
}

func NewRopeBuffer(contents []byte) *RopeBuffer {
	return &RopeBuffer{
		ropes.New(contents),
		nil,
		DetectLineDelim(contents),
	}
}

func (b *RopeBuffer) LineColToPos(line, col int) int {
	pos := b.getLineStartPos(line)

	// Have to do this algorithm for safety. If this function was declared to panic
	// or index out of bounds memory, if col > the given line length, it would be
	// more efficient and simpler. But unfortunately, I believe it is necessary.
	if col > 0 {
		_, r := b.rope.SplitAt(pos)
		l, _ := r.SplitAt(b.rope.Len() - pos)

		l.EachLeaf(func(n *ropes.Node) bool {
			data := n.Value() // Reference; not a copy.
			for _, r := range string(data) {
				if col == 0 || r == '\n' {
					return true // Found the position of the column
				}
				pos++
				col--
			}
			return false // Have not gotten to the appropriate position, yet
		})
	}

	return pos
}

func (b *RopeBuffer) Line(line int, delim bool) (bytes []byte, hasDelim bool) {
	pos := b.getLineStartPos(line)
	bytesLen := 0

	_, r := b.rope.SplitAt(pos)
	l, _ := r.SplitAt(b.rope.Len() - pos)

	l.EachLeaf(func(n *ropes.Node) bool {
		data := n.Value() // Reference; not a copy.
		delimPos := bytealg.Index(data, []byte(b.lineDelim))
		if delimPos != -1 { // Delim found
			if delim {
				delimPos = len(data) // Tells us to capture all of data
				hasDelim = true
			}
			bytesLen += delimPos
			return true
		}
		bytesLen += len(data)
		return false // Have not read the whole line, yet
	})

	return b.rope.Slice(pos, pos+bytesLen), hasDelim // NOTE: may be faster to do it ourselves
}

func (b *RopeBuffer) Slice(startLine, startCol, endLine, endCol int) []byte {
	endPos := b.LineColToPos(endLine, endCol)
	if length := b.rope.Len(); endPos >= length {
		endPos = length - 1
	}
	return b.rope.Slice(b.LineColToPos(startLine, startCol), endPos+1)
}

func (b *RopeBuffer) RuneAtPos(pos int) (val rune) {
	_, r := b.rope.SplitAt(pos)
	l, _ := r.SplitAt(b.rope.Len() - pos)

	l.EachLeaf(func(n *ropes.Node) bool {
		data := n.Value() // Reference; not a copy.
		val, _ = utf8.DecodeRune(data[0:])
		return true
	})

	return 0
}

func (b *RopeBuffer) EachRuneFromPos(pos int, f func(pos int, r rune) bool) {
	_, r := b.rope.SplitAt(pos)
	l, _ := r.SplitAt(b.rope.Len() - pos)

	l.EachLeaf(func(n *ropes.Node) bool {
		data := n.Value() // Reference; not a copy.
		for i, r := range string(data) {
			if f(pos+i, r) {
				return true
			}
		}
		return false
	})
}

func (b *RopeBuffer) Bytes() []byte {
	return b.rope.Value()
}

func (b *RopeBuffer) Insert(line, col int, value []byte) {
	b.rope.Insert(b.LineColToPos(line, col), value)
	b.shiftAnchors(line, col, utf8.RuneCount(value))
}

func (b *RopeBuffer) Remove(startLine, startCol, endLine, endCol int) {
	start := b.LineColToPos(startLine, startCol)
	end := b.LineColToPos(endLine, endCol) + 1

	if len := b.rope.Len(); end >= len {
		end = len
		if start > end {
			return
		}
	}

	b.rope.Remove(start, end)
	// Shift anchors within the range
	b.shiftAnchorsRemovedRange(start, startLine, startCol, endLine, endCol)
	// Shift anchors after the range
	b.shiftAnchors(endLine, endCol+1, start-end-1)
}

func (b *RopeBuffer) Count(startLine, startCol, endLine, endCol int, sequence []byte) int {
	startPos := b.LineColToPos(startLine, startCol)
	endPos := b.LineColToPos(endLine, endCol)
	return b.rope.Count(startPos, endPos+1, sequence)
}

func (b *RopeBuffer) Len() int {
	return b.rope.Len()
}

func (b *RopeBuffer) Lines() int {
	rope := b.rope
	return rope.Count(0, rope.Len(), []byte{'\n'}) + 1
}

func (b *RopeBuffer) LineDelimiter() string {
	return b.lineDelim
}

func (b *RopeBuffer) SetLineDelimiter(delim string) {
	b.lineDelim = delim
}

func (b *RopeBuffer) LineHasDelimiter(line int) bool {
	_, hasDelim := b.Line(line, true)
	return hasDelim // TODO: make this obsolete by providing delim bool params to functions
}

// getLineStartPos returns the first byte index of the given line (starting from zero).
// The returned index can be equal to the length of the buffer, not pointing to any byte,
// which means the byte is on the last, and empty, line of the buffer. If line is greater
// than or equal to the number of lines in the buffer, a panic is issued.
func (b *RopeBuffer) getLineStartPos(line int) int {
	var pos int

	if line > 0 {
		b.rope.IndexAllFunc(0, b.rope.Len(), []byte{'\n'}, func(idx int) bool {
			line--
			pos = idx + 1    // idx+1 = start of line after delimiter
			return line <= 0 // If pos is now the start of the line we're searching for
		})
	}

	if line > 0 { // If there aren't enough lines to reach line...
		panic("not enough lines in buffer to reach position")
	}

	return pos
}

func (b *RopeBuffer) RunesInLine(line int, delim bool) (runes int, hasDelim bool) {
	linePos := b.getLineStartPos(line)

	ropeLen := b.rope.Len()

	if linePos >= ropeLen {
		return 0, false
	}

	_, r := b.rope.SplitAt(linePos)
	l, _ := r.SplitAt(ropeLen - linePos)

	l.EachLeaf(func(n *ropes.Node) bool {
		data := n.Value() // Reference; not a copy.
		delimPos := bytealg.Index(data, []byte(b.lineDelim))
		if delimPos != -1 { // Delim found
			if delim {
				delimPos = len(data) // Causes us to capture all of data's length
				hasDelim = true
			}
			runes += utf8.RuneCount(data[:delimPos])
			return true
		}
		runes += utf8.RuneCount(data)
		return false // Have not read the whole line, yet
	})

	return
}

// ClampLineCol is a utility function to clamp any provided line and col to
// only possible values within the buffer, pointing to runes. It first clamps
// the line, then clamps the column. The column is clamped between zero and
// the last rune before the line delimiter.
func (b *RopeBuffer) ClampLineCol(line, col int) (int, int) {
	if line < 0 {
		line = 0
	} else if lines := b.Lines() - 1; line > lines {
		line = lines
	}

	if col < 0 {
		col = 0
	} else if runes, _ := b.RunesInLine(line, false); col > runes {
		col = runes
	}

	return line, col
}

// PosToLineCol converts a byte offset (position) of the buffer's bytes, into
// a line and column. Unless you are working with the Bytes() function, this
// is unlikely to be useful to you. Position will be clamped.
func (b *RopeBuffer) PosToLineCol(pos int) (int, int) {
	var line, col int
	var wasAtNewline bool

	if pos <= 0 {
		return line, col
	}

	b.rope.EachLeaf(func(n *ropes.Node) bool {
		data := n.Value()
		var i int
		for i < len(data) {
			if wasAtNewline { // Start of line
				if data[i] != '\n' { // If the start of this line does not happen to be a delim...
					wasAtNewline = false // Say we weren't previously at a delimiter
				}
				line, col = line+1, 0
			} else if data[i] == '\n' { // End of line
				wasAtNewline = true
				col++
			} else {
				col++
			}

			_, size := utf8.DecodeRune(data[i:])
			i += size
			pos -= size

			if pos < 0 {
				return true
			}
		}
		return false
	})

	return line, col
}

func (b *RopeBuffer) WriteTo(w io.Writer) (int64, error) {
	return b.rope.WriteTo(w)
}

// Currently meant for the Remove function: imagine if the removed region passes through a Cursor position.
// We want to shift the cursor to the start of the region, based upon where the cursor position is.
func (b *RopeBuffer) shiftAnchorsRemovedRange(startPos, startLine, startCol, endLine, endCol int) {
	for i := 0; i < len(b.anchors); i++ {
		v := b.anchors[i]
		if v == nil {
			b.removeCursorAtIdx(i)
			i--
			continue
		}

		if v.Line >= startLine && v.Line <= endLine {
			// If the anchor is not within the start or end columns
			if (v.Line == startLine && v.Col < startCol) || (v.Line == endLine && v.Col > endCol) {
				continue
			}
			cursorPos := b.LineColToPos(v.Line, v.Col)
			v.Line, v.Col = b.PosToLineCol(cursorPos + (startPos - cursorPos))
		}
	}
}

func (b *RopeBuffer) shiftAnchors(insertLine, insertCol, runeCount int) {
	for i := 0; i < len(b.anchors); i++ {
		v := b.anchors[i]
		if v == nil {
			b.removeCursorAtIdx(i)
			i--
			continue
		}
		if insertLine < v.Line || (insertLine == v.Line && insertCol <= v.Col) {
			v.Line, v.Col = b.PosToLineCol(b.LineColToPos(v.Line, v.Col) + runeCount)
		}
	}
}

// RegisterCursor adds the Cursor to a slice which the Buffer uses to update
// each Cursor based on changes that occur in the Buffer. Various functions are
// called on the Cursor depending upon where the edits occurred and how it should
// modify the Cursor's position. Unregister a Cursor before deleting it from
// memory, or forgetting it, with UnregisterPosition.
func (b *RopeBuffer) RegisterCursor(cursor *Cursor) {
	if cursor == nil {
		return
	}
	b.anchors = append(b.anchors, cursor)
}

// UnregisterCursor will remove the cursor from the list of watched Cursors.
// It is mandatory that a Cursor be unregistered before being freed from memory,
// or otherwise being forgotten.
func (b *RopeBuffer) UnregisterCursor(cursor *Cursor) {
	for i, v := range b.anchors {
		if cursor == v {
			b.removeCursorAtIdx(i)
			return
		}
	}
}

func (b *RopeBuffer) removeCursorAtIdx(idx int) {
	// Delete item at idx without preserving order
	b.anchors[idx] = b.anchors[len(b.anchors)-1]
	b.anchors[len(b.anchors)-1] = nil
	b.anchors = b.anchors[:len(b.anchors)-1]
}
