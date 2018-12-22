package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (ctrl *Controller) Run() error {
	ctrl.g.Highlight = true
	ctrl.g.Cursor = false
	ctrl.g.SelFgColor = gocui.ColorGreen

	ctrl.g.SetManagerFunc(layout)
	err := ctrl.bindKey()
	if err != nil {
		return err
	}

	initTitle(ctrl.g)

	return nil
}

// initTitle init sidebar title
func initTitle(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(sideViewName)
		if err != nil {
			log.Error("get view failed, ", err)
			return err
		}
		fmt.Fprintf(v, sidebarTitle())
		return nil
	})

}
