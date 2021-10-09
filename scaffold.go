package dos

import "github.com/gdamore/tcell/v2"

type Scaffold struct {
	MenuBar    *MenuBar
	MainWidget Widget
	Floating   []Widget
	focusIdx   int
}

func (s *Scaffold) IsMenuBarFocused() bool {
	return s.focusIdx == 0
}

func (s *Scaffold) IsMainWidgetFocused() bool {
	return s.focusIdx == 1
}

func (s *Scaffold) IsFloatingFocused() bool {
	return s.focusIdx == 2
}

func (s *Scaffold) FocusMenuBar() {
	s.setFocusMainWidget(false)
	s.setFocusFloating(false)
	s.setFocusMenuBar(true)
	s.focusIdx = 0
}

func (s *Scaffold) FocusMainWidget() {
	s.setFocusMenuBar(false)
	s.setFocusFloating(false)
	s.setFocusMainWidget(true)
	s.focusIdx = 1
}

func (s *Scaffold) FocusFloating() {
	s.setFocusMenuBar(false)
	s.setFocusMainWidget(false)
	s.setFocusFloating(true)
	s.focusIdx = 2
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

func (s *Scaffold) setFocusFloating(v bool) {
	if len(s.Floating) > 0 {
		s.Floating[len(s.Floating)-1].SetFocused(v)
	}
}

func (s *Scaffold) mainWidgetRect(currentRect Rect) Rect {
	rect := currentRect
	if s.MenuBar != nil {
		_, h := s.MenuBar.DisplaySize(currentRect.W, currentRect.H)
		rect.Y += h
		rect.H -= h
	}
	w, h := s.MainWidget.DisplaySize(rect.W, rect.H)
	return Rect{rect.X, rect.Y, w, h}
}

func (s *Scaffold) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if len(s.Floating) > 0 {
		for i := len(s.Floating) - 1; i >= 0; i-- {
			if s.Floating[i].HandleMouse(currentRect, ev) {
				// This window handled the event, if they are not at the end of
				// the slice, then we move them there (so they're drawn last).
				if i != len(s.Floating)-1 {
					s.setFocusFloating(false) // Unfocus current top window
					win := s.Floating[i]      // Make a copy of item at i
					// Shift items after i left
					copy(s.Floating[i:], s.Floating[i+1:])
					s.Floating[len(s.Floating)-1] = win // Move copy to end
				}

				s.setFocusMenuBar(false)
				s.setFocusMainWidget(false)
				s.focusIdx = 2
				return true
			}
		}
	}
	if s.MainWidget != nil {
		if s.MainWidget.HandleMouse(s.mainWidgetRect(currentRect), ev) {
			s.setFocusMenuBar(false)
			s.setFocusFloating(false)
			s.focusIdx = 1
			return true
		}
	}
	if s.MenuBar != nil {
		sizeX, sizeY := s.MenuBar.DisplaySize(currentRect.W, currentRect.H)
		if s.MenuBar.HandleMouse(Rect{currentRect.X, currentRect.Y, sizeX, sizeY}, ev) {
			s.setFocusMainWidget(false)
			s.setFocusFloating(false)
			s.focusIdx = 0
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
		s.focusIdx = 2
	} else if s.MainWidget != nil {
		s.MainWidget.SetFocused(b)
		s.focusIdx = 1
	} else if s.MenuBar != nil {
		s.MenuBar.SetFocused(b)
		s.focusIdx = 0
	}
}

func (s *Scaffold) DisplaySize(boundsW, boundsH int) (w, h int) {
	return boundsW, boundsH
}

func (s *Scaffold) Draw(rect Rect, screen tcell.Screen) {
	if s.MainWidget != nil {
		s.MainWidget.Draw(s.mainWidgetRect(rect), screen)
	}
	if s.MenuBar != nil {
		s.MenuBar.Draw(Rect{0, 0, rect.W, 1}, screen)
	}
	for i := 0; i < len(s.Floating); i++ { // Draw back to front
		s.Floating[i].Draw(rect, screen)
	}
}
