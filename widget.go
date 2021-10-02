package dos

import (
	"github.com/gdamore/tcell/v2"
)

type Widget interface {
	// HandleMouse is called by a parent of the Widget when they receive the
	// event and have found that the event occurred within the bounds of their
	// Widget. The function will always check if the event occurred within the
	// bounds of the Widget before continuing. HandleMouse will return `true` if
	// the event is handled. Otherwise, `false`, so the event continues to be
	// propagated. If this Widget or any of its child Widgets handle this event
	// successfully (by returning true), then SetFocused(true) should be called.
	//
	// The currentRect is used by the implementor to determine its current
	// position and size. The size is determined its parent.
	HandleMouse(currentRect Rect, ev *tcell.EventMouse) bool
	// HandleKey is called by a parent of the Widget when they receive the
	// event. The Widget will only try to handle the event if it is focused.
	// HandleKey will return `true` if the event is handled. Otherwise, `false`,
	// so the event can continue to be propagated. If this Widget or any of its
	// child Widgets handle this event successfully (by returning true), then
	// SetFocused(true) should be called.
	//
	// The currentRect is used by the implementor to determine its current
	// position and size. The size is determined its parent.
	HandleKey(ev *tcell.EventKey) bool
	// SetFocused alerts the Widget that it has received input focus from the
	// user. The value can be kept in the Widget to differ its appearance during
	// Draw. The Widget will call SetFocused(b) on all of its children, also.
	SetFocused(b bool)
	// DisplaySizeInBounds returns the expected size of the Widget when it will
	// be drawn. This is primarily used for containers like Center that require
	// the Widget's size when determining how to center it.
	DisplaySize(boundsW, boundsH int) (w, h int)
	// Draw renders the Widget onto the terminal screen, bounded by the provided
	// Rect. It is a bug if the Widget draws any part of itself outside the rect
	// provided. Draw should not call Sync() on the tcell.Screen or other
	// synchronizing functions, as all synchronization will be done by the event
	// loop.
	Draw(rect Rect, s tcell.Screen)
}
