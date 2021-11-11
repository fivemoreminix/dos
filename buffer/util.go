package buffer

import (
	"bytes"
)

const (
	LF   = "\n"
	CRLF = "\r\n"
)

// Max returns the larger integer.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller integer.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Clamp keeps `v` within `a` and `b` numerically. `a` must be smaller than `b`.
// Returns clamped `v`.
func Clamp(v, a, b int) int {
	return Max(a, Min(v, b))
}

// DetectLineDelim searches for a CRLF "\r\n" or LF "\n", and if neither is found,
// it produces a default value LF "\n".
func DetectLineDelim(contents []byte) string {
	lfpos := bytes.IndexByte(contents, '\n')
	if lfpos <= 0 {
		return LF
	}
	if contents[lfpos-1] == '\r' {
		return CRLF
	}
	return LF
}
