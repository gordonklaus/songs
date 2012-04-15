package gui

import "image"

type MouseHandler interface {
	MousePressed(button int, p image.Point)
	MouseDragged(button int, p image.Point)
	MouseReleased(button int, p image.Point)
}


type AggregateMouseHandler []MouseHandler
func (a AggregateMouseHandler) MousePressed(button int, p image.Point) {
	for _, h := range a { h.MousePressed(button, p) }
}
func (a AggregateMouseHandler) MouseDragged(button int, p image.Point) {
	for _, h := range a { h.MouseDragged(button, p) }
}
func (a AggregateMouseHandler) MouseReleased(button int, p image.Point) {
	for _, h := range a { h.MouseReleased(button, p) }
}


type ClickHandler struct {
	f func(int, image.Point)
}
func NewClickHandler(f func(int, image.Point)) *ClickHandler {
	return &ClickHandler{f}
}
func (c ClickHandler) MousePressed(button int, p image.Point) {
	c.f(button, p)
}
func (c ClickHandler) MouseDragged(button int, p image.Point) {}
func (c ClickHandler) MouseReleased(button int, p image.Point) {}


func NewClickKeyboardFocuser(view View) *ClickHandler {return NewClickHandler(func(int, image.Point) { view.TakeKeyboardFocus() }) }


type ViewDragger struct {
	view View
	p image.Point
}
func NewViewDragger(view View) *ViewDragger {
	return &ViewDragger{view:view}
}
func (d *ViewDragger) MousePressed(button int, p image.Point) {
	d.view.Raise()
	d.p = p
}
func drag(v View, p1, p2 image.Point) { v.Move(v.Position().Add(p2.Sub(p1))) }
func (d ViewDragger) MouseDragged(button int, p image.Point) { drag(d.view, d.p, p) }
func (d ViewDragger) MouseReleased(button int, p image.Point) { drag(d.view, d.p, p) }


type ViewPanner struct {
	view View
	p image.Point
}
func NewViewPanner(view View) *ViewPanner {
	return &ViewPanner{view:view}
}
func (vp *ViewPanner) MousePressed(button int, p image.Point) {
	vp.p = p
}
func pan(v View, p1, p2 image.Point) { v.Pan(p1.Sub(p2)) }
func (vp *ViewPanner) MouseDragged(button int, p image.Point) { pan(vp.view, vp.p, p) }
func (vp *ViewPanner) MouseReleased(button int, p image.Point) { pan(vp.view, vp.p, p) }
