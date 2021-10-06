package dos

import "github.com/gdamore/tcell/v2"

type Scaffold struct {
	MenuBar    *MenuBar
	MainWidget Widget
	Floating   []Widget
}

func (s *Scaffold) setFocusMenuBar(v bool) {
	if s.MenuBar != nil {
		s.MenuBar.SetFocused(v)
	}
}

func (s *Scaffold) setFocusMainWidget(v bool) {
	if s.MainWidget != nil {
		s.MainWidget.SetFocused(v)
	}
}

func (s *Scaffold) setFocusWindows(v bool) {
	if len(s.Floating) > 0 {
		s.Floating[len(s.Floating)-1].SetFocused(v)
	}
}

func (s *Scaffold) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if len(s.Floating) > 0 {
		for i := len(s.Floating) - 1; i >= 0; i-- {
			if s.Floating[i].HandleMouse(currentRect, ev) {
				// This window handled the event, if they are not at the end of
				// the slice, then we move them there (so they're drawn last).
				if i != len(s.Floating)-1 {
					s.setFocusWindows(false) // Unfocus current top window
					win := s.Floating[i]     // Make a copy of item at i
					// Shift items after i left
					copy(s.Floating[i:], s.Floating[i+1:])
					s.Floating[len(s.Floating)-1] = win // Move copy to end
				}

				s.setFocusMenuBar(false)
				s.setFocusMainWidget(false)
				return true
			}
		}
	}
	if s.MainWidget != nil {
		w, h := s.MainWidget.DisplaySize(currentRect.W, currentRect.H)
		if s.MainWidget.HandleMouse(Rect{currentRect.X, currentRect.Y, w, h}, ev) {
			s.setFocusMenuBar(false)
			s.setFocusWindows(false)
			return true
		}
	}
	if s.MenuBar != nil {
		sizeX, sizeY := s.MenuBar.DisplaySize(currentRect.W, currentRect.H)
		if s.MenuBar.HandleMouse(Rect{currentRect.X, currentRect.Y, sizeX, sizeY}, ev) {
			s.setFocusMainWidget(false)
			s.setFocusWindows(false)
			return true
		}
	}
	return false
}

func (s *Scaffold) HandleKey(ev *tcell.EventKey) bool {
	if len(s.Floating) > 0 {
		for i := len(s.Floating) - 1; i >= 0; i-- {
			if s.Floating[i].HandleKey(ev) {
				return true
			}
		}
	} else {
		if s.MainWidget != nil {
			if s.MainWidget.HandleKey(ev) {
				return true
			}
		}
		if s.MenuBar != nil {
			return s.MenuBar.HandleKey(ev)
		}
	}
	return false
}

func (s *Scaffold) SetFocused(b bool) {
	if len(s.Floating) > 0 {
		s.Floating[len(s.Floating)-1].SetFocused(b)
	} else if s.MainWidget != nil {
		s.MainWidget.SetFocused(b)
	} else if s.MenuBar != nil {
		s.MenuBar.SetFocused(b)
	}
}

func (s *Scaffold) DisplaySize(boundsW, boundsH int) (w, h int) {
	return boundsW, boundsH
}

func (s *Scaffold) Draw(rect Rect, screen tcell.Screen) {
	if s.MainWidget != nil {
		mainRect := rect
		if s.MenuBar != nil {
			// TODO: should use s.MenuBar.DisplaySize function instead
			mainRect.Y += 1
			mainRect.H -= 1
		}
		s.MainWidget.Draw(mainRect, screen)
	}
	if s.MenuBar != nil {
		s.MenuBar.Draw(Rect{0, 0, rect.W, 1}, screen)
	}
	for i := 0; i < len(s.Floating); i++ { // Draw back to front
		s.Floating[i].Draw(rect, screen)
	}
}
