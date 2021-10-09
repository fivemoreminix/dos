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
	Decorated      bool           // Whether to draw a styled box around the Menu
	Decoration     *BoxDecoration // Used only if Decorated is true
	NormalStyle    tcell.Style
	SelectionStyle tcell.Style
	Selected       int
	itemsExpanded  bool // True if the Selected submenu should be visible
}

func (m *Menu) ActivateItem(idx int) {
	switch item := m.Items[idx]; item.Type {
	case MenuItemAction:
		if item.Action != nil {
			item.Action()
		}
	case MenuItemSubmenu:
	case MenuItemSeparator:
	}
}

func (m *Menu) HandleMouse(rect Rect, ev *tcell.EventMouse) bool {
	if ev.Buttons() == tcell.ButtonPrimary {
		sizeW, sizeH := m.DisplaySize(0, 0)
		posX, posY := ev.Position()
		// Check if user clicked any part of the menu (including border)
		if posX >= rect.X && posX < rect.X+sizeW && posY >= rect.Y && posY < rect.Y+sizeH {
			offsetX, offsetY := 0, 0
			if m.Decorated { // If we have a border around the menu
				sizeW -= 2
				offsetX, offsetY = 1, 1
			}
			// Check that the click occurred between the borders of the menu
			if posX >= rect.X+offsetX && posX < rect.X+offsetX+sizeW {
				for i := 0; i < len(m.Items); i++ {
					if posY == rect.Y+i+offsetY {
						m.Selected = i
						m.ActivateItem(i)
					}
				}
			}
			return true // User clicked somewhere in the menu
		}
	}
	return false
}

func (m *Menu) HandleKey(ev *tcell.EventKey) bool {
	// TODO: support submenus in Menu
	if m.Items != nil && len(m.Items) > 0 {
		switch ev.Key() {
		case tcell.KeyUp:
			for {
				m.Selected--
				if m.Selected < 0 {
					m.Selected = len(m.Items) - 1
				}
				if m.Items[m.Selected].Type != MenuItemSeparator {
					break // TODO: Loop has the potential to block forever
				}
			}
		case tcell.KeyDown:
			for {
				m.Selected++
				if m.Selected >= len(m.Items) {
					m.Selected = 0
				}
				if m.Items[m.Selected].Type != MenuItemSeparator {
					break
				}
			}
		case tcell.KeyEnter:
			m.ActivateItem(m.Selected)
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
	if m.Decorated {
		return widestItemWidth + 2, len(m.Items) + 2
	} else {
		return widestItemWidth, len(m.Items)
	}
}

func (m *Menu) Draw(rect Rect, s tcell.Screen) {
	width, _ := m.DisplaySize(0, 0)
	offsetX, offsetY := 0, 0
	decoration := m.Decoration
	if m.Decorated {
		offsetX = 1
		offsetY = 1
		if decoration == nil {
			decoration = &DefaultBoxDecoration
		}
		DrawBox(rect, decoration, s)
	}
	for i := 0; i < len(m.Items); i++ {
		style := m.NormalStyle
		if i == m.Selected {
			style = m.SelectionStyle
		}

		if m.Items[i].Type == MenuItemSeparator && m.Decorated {
			s.SetContent(rect.X, rect.Y+offsetY+i, decoration.JointL, nil, decoration.Style)
			s.SetContent(rect.X+width-1, rect.Y+offsetY+i, decoration.JointR, nil, decoration.Style)
			for col := 0; col < width-offsetX*2; col++ {
				s.SetContent(rect.X+offsetX+col, rect.Y+offsetY+i, decoration.Hor, nil, decoration.Style)
			}
		} else {
			for col := 0; col < width-offsetX*2; col++ {
				s.SetContent(rect.X+offsetX+col, rect.Y+offsetY+i, ' ', nil, style)
			}
			DrawString(rect.X+offsetX, rect.Y+offsetY+i, m.Items[i].Title, style, s)
		}
	}
}
