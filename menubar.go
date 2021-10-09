package dos

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// A MenuBarItem is a Menu with an added Title field.
type MenuBarItem struct {
	Title string
	Menu
}

type MenuBar struct {
	Menus          []MenuBarItem
	NormalStyle    tcell.Style
	SelectionStyle tcell.Style
	Selected       int
	expanded       bool // Whether the user is expanding the menus currently
	focused        bool // Whether to accept keyboard input and highlight selection
}

// ItemRects returns a slice of Rects for each Menu's title that the user
// selects before expanding the actual Menu. The returned slice length will be
// equal to the length of Menus.
func (m *MenuBar) ItemRects(rect Rect) []Rect {
	// TODO: MenuBar ItemRects() may need to be optimized by removing the allocation
	rects := make([]Rect, len(m.Menus))
	col := 1
	for i := 0; i < len(m.Menus); i++ {
		textWidth := runewidth.StringWidth(m.Menus[i].Title)
		rects[i] = Rect{rect.X + col, rect.Y, textWidth + 2, 1}
		col += textWidth + 2
	}
	return rects
}

func (m *MenuBar) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	rects := m.ItemRects(currentRect)
	if ev.Buttons()&tcell.ButtonPrimary != 0 {
		for i := 0; i < len(rects); i++ {
			if rects[i].HasPoint(ev.Position()) {
				m.focused = true
				m.Selected = i
				m.expanded = true
				return true
			}
		}
	}
	if m.expanded && len(m.Menus) > 0 {
		rect := Rect{rects[m.Selected].X, rects[m.Selected].Y + 1, currentRect.W, currentRect.H}
		m.Menus[m.Selected].HandleMouse(rect, ev)
	}
	return false
}

func (m *MenuBar) HandleKey(ev *tcell.EventKey) bool {
	if m.focused && len(m.Menus) > 0 {
		if m.expanded {
			if m.Menus[m.Selected].HandleKey(ev) {
				return true
			}
		}

		switch ev.Key() {
		case tcell.KeyLeft:
			// Reset menu selection idx before changing
			m.Menus[m.Selected].Selected = 0
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.Menus) - 1
			}
		case tcell.KeyRight:
			// Reset menu selection before changing
			m.Menus[m.Selected].Selected = 0
			m.Selected++
			if m.Selected >= len(m.Menus) {
				m.Selected = 0
			}
		case tcell.KeyEnter:
			m.expanded = !m.expanded
		default:
			return false
		}
		return true
	}
	return false
}

func (m *MenuBar) SetFocused(b bool) {
	m.focused = b
	if !b {
		m.expanded = false
	}
	m.Menus[m.Selected].Selected = 0
	// NOTE: I am not calling SetFocused on the highlighted menu because currently
	// menus do not accept focus.
}

func (m *MenuBar) DisplaySize(boundsW, _ int) (w, h int) {
	return boundsW, 1
}

func (m *MenuBar) Draw(rect Rect, s tcell.Screen) {
	for col := 0; col < rect.W; col++ {
		s.SetContent(rect.X+col, rect.Y, ' ', nil, m.NormalStyle)
	}
	if len(m.Menus) > 0 {
		rects := m.ItemRects(rect)
		for i, r := range rects {
			style := m.NormalStyle
			if m.focused && i == m.Selected { // If this menu is selected
				style = m.SelectionStyle

				if m.expanded { // If the selected menu is also expanded
					menuW, menuH := m.Menus[i].DisplaySize(0, 0)
					m.Menus[i].Draw(Rect{r.X, r.Y + 1, menuW, menuH}, s)
				}
			}
			DrawString(r.X, r.Y, fmt.Sprintf(" %s ", m.Menus[i].Title), style, s)
		}
	}
}
