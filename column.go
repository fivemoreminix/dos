package dos

import "github.com/gdamore/tcell/v2"

// A Column orders its children vertically.
type Column struct {
	Children []Widget
	Align    Align
}

func (c Column) HandleClick(ev *tcell.EventMouse) bool {
	panic("implement me")
}

func (c Column) HandleKey(ev *tcell.EventKey) bool {
	panic("implement me")
}

func (c Column) SetFocused(b bool) {
	panic("implement me")
}

func (c Column) DisplaySizeInBounds(boundsW, boundsH int) (w, h int) {
	panic("implement me")
}

func (c Column) Draw(rect Rect, s tcell.Screen) {
	panic("implement me")
}
