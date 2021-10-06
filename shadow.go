package dos

import "github.com/gdamore/tcell/v2"

// A Shadow draws a shadow covering one cell below, and two cells right of its
// bounding box. This is useful for dialog windows or buttons that need depth.
type Shadow struct {
	Child     Widget
	Style     tcell.Style
	MakeSmall bool
}

func (s *Shadow) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if s.Child != nil {
		return s.Child.HandleMouse(currentRect, ev)
	}
	return false
}

func (s *Shadow) HandleKey(ev *tcell.EventKey) bool {
	if s.Child != nil {
		return s.Child.HandleKey(ev)
	}
	return false
}

func (s *Shadow) SetFocused(b bool) {
	if s.Child != nil {
		s.Child.SetFocused(b)
	}
}

func (s *Shadow) DisplaySize(boundsW, boundsH int) (w, h int) {
	if s.Child != nil {
		return s.Child.DisplaySize(boundsW, boundsH)
	}
	return 0, 0
}

func (s *Shadow) shadowCell(x, y int, screen tcell.Screen) (width int) {
	r, combc, _, width := screen.GetContent(x, y)
	screen.SetContent(x, y, r, combc, s.Style)
	return width
}

// Draw causes the Shadow to intentionally set cells outside its provided rect.
// The provided rect is passed directly to the child.
func (s *Shadow) Draw(rect Rect, screen tcell.Screen) {
	// TODO: make Shadow extents configurable

	reversedStyle := s.Style.Reverse(true)

	// Right side
	width := 2
	if s.MakeSmall {
		// Draw this in the top right corner
		screen.SetContent(rect.X+rect.W, rect.Y, '▄', nil, reversedStyle)
		width = 1
	}
	for row := rect.Y + 1; row < rect.Y+rect.H; row++ {
		// Draw side two columns wide
		for col := rect.X + rect.W; col < rect.X+rect.W+width; col++ {
			width := s.shadowCell(col, row, screen)
			if width > 1 {
				// If we are in the first iteration of the col loop, this will
				// prevent us from accessing the next col of an east-asian rune.
				col++
			}
		}
	}
	// Bottom side
	for col := rect.X + 1; col < rect.X+rect.W+width; col++ {
		if s.MakeSmall {
			screen.SetContent(col, rect.Y+rect.H, '▀', nil, reversedStyle)
		} else {
			width := s.shadowCell(col, rect.Y+rect.H, screen)
			if width > 1 {
				col++ // Step over additional cell of east-asian characters
			}
		}
	}
	// Draw child
	if s.Child != nil {
		s.Child.Draw(rect, screen)
	}
}
