package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	modalMaxWidth  = 56
	modalMaxHeight = 20

	compactModalWidth  = 48
	compactModalHeight = 10
)

type fittedCenter struct {
	*tview.Box
	child     tview.Primitive
	maxWidth  int
	maxHeight int
}

func newFittedCenter(child tview.Primitive) *fittedCenter {
	return newSizedCenter(child, modalMaxWidth, modalMaxHeight)
}

func newCompactCenter(child tview.Primitive) *fittedCenter {
	return newSizedCenter(child, compactModalWidth, compactModalHeight)
}

func newSizedCenter(child tview.Primitive, maxWidth, maxHeight int) *fittedCenter {
	return &fittedCenter{
		Box:       tview.NewBox(),
		child:     child,
		maxWidth:  maxWidth,
		maxHeight: maxHeight,
	}
}

func (center *fittedCenter) Draw(screen tcell.Screen) {
	x, y, width, height := center.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}
	childWidth := min(center.maxWidth, max(width-2, 1))
	childHeight := min(center.maxHeight, max(height-2, 1))
	childX := x + (width-childWidth)/2
	childY := y + (height-childHeight)/2
	center.child.SetRect(childX, childY, childWidth, childHeight)
	center.child.Draw(screen)
}

func (center *fittedCenter) Focus(delegate func(p tview.Primitive)) {
	delegate(center.child)
}

func (center *fittedCenter) HasFocus() bool {
	return center.child.HasFocus()
}

func (center *fittedCenter) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return center.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if handler := center.child.InputHandler(); handler != nil {
			handler(event, setFocus)
		}
	})
}

func (center *fittedCenter) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return center.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if handler := center.child.MouseHandler(); handler != nil {
			return handler(action, event, setFocus)
		}
		return false, nil
	})
}
