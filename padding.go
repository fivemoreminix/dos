package dos

import "github.com/gdamore/tcell/v2"

type Padding struct {
	Child  Widget
	Top    int
	Right  int
	Bottom int
	Left   int
}

func (p *Padding) GetChildRect(currentRect Rect) Rect {
	return Rect{
		currentRect.X + p.Left,
		currentRect.Y + p.Top,
		currentRect.W - p.Left - p.Right,
		currentRect.H - p.Top - p.Bottom,
	}
}

func (p *Padding) HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool {
	if p.Child != nil {
		return p.Child.HandleMouse(p.GetChildRect(currentRect), ev)
	}
	return false
}

func (p *Padding) HandleKey(ev *tcell.EventKey) bool {
	if p.Child != nil {
		return p.Child.HandleKey(ev)
	}
	return false
}

func (p *Padding) SetFocused(b bool) {
	if p.Child != nil {
		p.Child.SetFocused(b)
	}
}

func (p *Padding) DisplaySize(boundsW, boundsH int) (w, h int) {
	if p.Child != nil {
		w, h = p.Child.DisplaySize(boundsW-p.Left-p.Right, boundsH-p.Top-p.Bottom)
		return w + p.Left + p.Right, h + p.Top + p.Bottom
	}
	return 0, 0
}

func (p *Padding) Draw(rect Rect, s tcell.Screen) {
	if p.Child != nil {
		DrawRect(p.GetChildRect(rect), ' ', tcell.Style{}.Background(tcell.ColorYellow), s)
		p.Child.Draw(p.GetChildRect(rect), s)
	}
}
