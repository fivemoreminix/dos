package dos

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type MenuItemType uint8

const (
	MenuItemAction MenuItemType = iota
	MenuItemSubmenu
	MenuItemSeparator
)

// A MenuItem is a selectable option inside a Menu. If the Type is
// MenuItemAction, then only the Action field should be accessed. Likewise, if
// the Type is MenuItemSubmenu, then only the Submenu field should be accessed.
type MenuItem struct {
	Title   string
	Type    MenuItemType
	Action  func()
	Submenu *Menu
}

// A Menu contains a list of selectable items the user may either click or
// scroll through using the arrow keys. The items of a Menu are the type
// MenuItem, which can be an action or expand into another menu, known as a
// submenu. This widget will handle any events it receives, so do not pass
// an event to the Menu if it is not focused or otherwise visible.
type Menu struct {
	Items          []MenuItem
	Decorated      bool // Whether to draw a styled box around the Menu
	NormalStyle    tcell.Style
	SelectionStyle tcell.Style
	Selected       int
	itemsExpanded  bool // True if the Selected submenu should be visible
}

func (m *Menu) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	return false // TODO: mouse support for Menu
}

func (m *Menu) HandleKey(ev *tcell.EventKey) bool {
	// TODO: support submenus in Menu
	if m.Items != nil && len(m.Items) > 0 {
		switch ev.Key() {
		case tcell.KeyUp:
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.Items) - 1
			}
		case tcell.KeyDown:
			m.Selected++
			if m.Selected >= len(m.Items) {
				m.Selected = 0
			}
		case tcell.KeyEnter:
			switch item := m.Items[m.Selected]; item.Type {
			case MenuItemAction:
				if item.Action != nil {
					item.Action()
				}
			case MenuItemSubmenu:
			case MenuItemSeparator:
			}
		default:
			return false
		}
		return true
	}
	return false
}

func (m *Menu) SetFocused(bool) {}

// DisplaySize for Menu returns the minimum size to show the Menu normally,
// so it ignores the input boundary values.
func (m *Menu) DisplaySize(int, int) (w int, h int) {
	// Find the widest item
	widestItemWidth := 0
	for i := range m.Items {
		if width := runewidth.StringWidth(m.Items[i].Title); width > widestItemWidth {
			widestItemWidth = width
		}
	}
	return widestItemWidth, len(m.Items)
}

func (m *Menu) Draw(rect Rect, s tcell.Screen) {
	width, _ := m.DisplaySize(0, 0)
	for row := 0; row < len(m.Items); row++ {
		style := m.NormalStyle
		if row == m.Selected {
			style = m.SelectionStyle
		}

		for col := 0; col < width; col++ {
			s.SetContent(rect.X+col, rect.Y+row, ' ', nil, style)
		}
		DrawString(rect.X, rect.Y+row, m.Items[row].Title, style, s)
	}
}
