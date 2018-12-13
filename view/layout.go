package view

import (
	"github.com/jroimartin/gocui"
)

var mainViewName = "main"
var sideViewName = "side"
var cmdlineViewName = "cmdline"

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(sideViewName, 0, 0, int(0.2*float32(maxX)), maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			log.Error(err)
			return err
		}
		v.Wrap = false
		v.Frame = true
		v.Autoscroll = true
		v.Highlight = true
		v.Overwrite = true

		if _, err := setCurrentViewOnTop(g, sideViewName); err != nil {
			log.Error(err)
			return err
		}
	}

	if v, err := g.SetView(mainViewName, int(0.2*float32(maxX)), 0, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			log.Error(err)
			return err
		}
		v.Wrap = true
		v.Frame = true
		v.Autoscroll = true
	}

	if v, err := g.SetView(cmdlineViewName, 0, maxY-2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			log.Error(err)
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
	}

	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		log.Error(err)
		return nil, err
	}
	return g.SetViewOnTop(name)
}
