package view

import (
	ui "github.com/aiportal/termui"
)

type GridView struct {
	controls []*GridViewControl
	rows     []*ui.Row
}

func (v *GridView) Show() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// build layout
	ui.Body.AddRows(v.rows...)

	ui.Body.Align()
	ui.Clear()
	ui.Render(ui.Body)

	for _, c := range v.controls {
		if c.Handlers == nil {
			continue
		}
		for k, v := range c.Handlers {
			ui.Handle(k, v)
		}
	}

	ui.Handle("sys/wnd/resize", func(e ui.Event) {
		evt := e.Data.(ui.EvtWnd)
		ui.Body.Width = evt.Width
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}

func (v *GridView) Update() {
	ui.Body.Align()
	ui.Clear()
	ui.Render(ui.Body)
}

func (v *GridView) Control(name string) *GridViewControl {
	for _, c := range v.controls {
		if c.Name == name {
			return c
		}
	}
	return nil
}

type HandlerMap map[string]func(ui.Event)

type GridViewControl struct {
	Name     string
	Element  interface{}
	Handlers HandlerMap
	Data 	 interface{}
}

func NewGridViewControl(name string, element ui.GridBufferer, handlers HandlerMap) *GridViewControl {
	c := new(GridViewControl)
	c.Name = name
	c.Element = element
	c.Handlers = handlers
	return c
}
