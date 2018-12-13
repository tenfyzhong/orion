package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Run 运行
func Run() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g.Highlight = true
	g.Cursor = false
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)
	err = bindKey(g)
	if err != nil {
		return nil, err
	}

	initTitle(g)

	return g, nil
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
