package buffer

import (
	"io"
)

// A Buffer is wrapper around any buffer data structure like a rope or gap buffer
// that can be used for editing text. One way this interface helps is by making
// all API function parameters line and column indexes, so it is simple and easy
// to index and use like a text editor. All lines and columns start at zero, and
// all "end" ranges are inclusive.
//
// Any bounds out of range are panics! If you are unsure your position or range
// may be out of bounds, use ClampLineCol() or compare with Lines() or ColsInLine().
type Buffer interface {
	// Line gets a slice of the provided line with the delimiter if delim is true,
	// and one is present. Returns line bytes and whether a delimiter is included in
	// the result. Data returned may or may not be a copy: do not write to it.
	Line(line int, delim bool) (bytes []byte, hasDelim bool)

	// Returns a slice of the buffer from startLine, startCol, to endLine, endCol,
	// inclusive bounds. The returned value may or may not be a copy of the data,
	// so do not write to it.
	Slice(startLine, startCol, endLine, endCol int) []byte

	// RuneAtPos returns the UTF-8 rune at the byte position `pos` of the buffer. The
	// position must be a correct position, otherwise zero is returned.
	RuneAtPos(pos int) rune

	// EachRuneFromPos executes the function `f` for each rune from byte position `pos`.
	// This function should be used as opposed to performing a "per character" operation
	// manually, as it enables caching buffer operations and safety checks. The function
	// returns when the end of the buffer is met or `f` returns true.
	EachRuneFromPos(pos int, f func(pos int, r rune) bool)

	// Bytes returns all of the bytes in the buffer. This function is very likely
	// to copy all of the data in the buffer. Use sparingly. Try using other methods,
	// where possible.
	Bytes() []byte

	// Insert copies a byte slice (inserting it) into the position at line, col.
	Insert(line, col int, bytes []byte)

	// Remove deletes any characters between startLine, startCol, and endLine,
	// endCol, inclusive bounds.
	Remove(startLine, startCol, endLine, endCol int)

	// Returns the number of occurrences of 'sequence' in the buffer, within the range
	// of start line and col, to end line and col, inclusive bounds.
	Count(startLine, startCol, endLine, endCol int, sequence []byte) int

	// Len counts the number of bytes in the buffer.
	Len() int

	// Lines counts the number of lines in the buffer. If the buffer is empty, 1 is
	// returned, because there is always at least one line. This function basically
	// counts 1 + the number of delimiters in a buffer.
	Lines() int

	// Returns the line delimiter being used by the Buffer to separate lines. A Buffer
	// should have their delimiter set automatically, most likely when provided text.
	LineDelimiter() string

	// SetLineDelimiter overwrites the current delimiter being used by the Buffer to
	// separate the contents into lines.
	SetLineDelimiter(delim string)

	// LineHasDelimiter returns true if the line ends with the value of the Buffer's
	// set line delimiter. You can get the line delimiter string with LineDelimiter(),
	// and change it with SetLineDelimiter(delim).
	LineHasDelimiter(line int) bool

	// RunesInLine returns the number of runes in the given line. That is, the
	// number of UTF-8 codepoints in the line, not bytes. Includes the line delimiter
	// in the count. If that line delimiter is CRLF ('\r\n'), then it adds two.
	//RunesInLineWithDelim(line int) int

	// RunesInLine returns the number of runes in the given line. That is, the
	// number of UTF-8 codepoints in the line, not bytes. delim can be true to include
	// the line delimiter if present. Returns the number of runes and whether a
	// line delimiter is included in the result.
	RunesInLine(line int, delim bool) (runes int, hasDelim bool)

	// ClampLineCol is a utility function to clamp any provided line and col to
	// only possible values within the buffer pointing to runes. It first clamps
	// the line, then clamps the column. The column is clamped between zero and
	// the last rune before the line delimiter if one is present.
	ClampLineCol(line, col int) (int, int)

	// LineColToPos returns the index of the byte at line, col. If col is greater
	// than the length of the line, the position of the last byte in the line is
	// returned, instead. May include a delimiter. Use LineHasDelimiter(line) to
	// check if one is present.
	LineColToPos(line, col int) int

	// PosToLineCol converts a byte offset (position) of the buffer's bytes, into
	// a line and column. Position will be clamped.
	PosToLineCol(pos int) (line, col int)

	// Writes the Buffer to the provided io.Writer. Returns the number of bytes
	// written, and any error that may have occurred.
	WriteTo(w io.Writer) (int64, error)

	// RegisterCursor adds the Cursor to a slice which the Buffer manages to update
	// each Cursor based on changes that occur in the Buffer. Various functions are
	// called on the Cursor depending upon where the edits occurred and how it should
	// modify the Cursor's position.
	//
	// It is a good idea to call UnregisterCursor() when the cursor is no longer needed,
	// otherwise it will persist.
	RegisterCursor(cursor *Cursor)

	// UnregisterCursor will remove the cursor from the list of watched Cursors.
	UnregisterCursor(cursor *Cursor)
}
