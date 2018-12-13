package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/tenfyzhong/orion/controller"
	"github.com/tenfyzhong/orion/model"
)

// NewUpdateFunc get a update function
func NewUpdateFunc(g *gocui.Gui) controller.UpdateFunc {
	return func(m *model.Message) {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View(sideViewName)
			if err != nil {
				log.Error("get view failed", err)
				return err
			}
			output := messageSidebarString(m)
			fmt.Fprintf(v, "\n%s", output)
			// set the cursor to the last line
			_, y := v.Cursor()
			_, maxY := v.Size()
			newY := y + 1
			if newY >= maxY-1 {
				newY = maxY - 1
			}
			v.SetCursor(0, newY)
			return nil
		})
	}
}
