package view

import (
	"errors"

	"github.com/jroimartin/gocui"
)

func bindKey(g *gocui.Gui) error {
	if g == nil {
		return errors.New("gui is nil")
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Critical(err)
	}

	if err := g.SetKeybinding(mainViewName, gocui.KeyArrowLeft, gocui.ModNone, setSideOnTop); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, gocui.KeyArrowRight, gocui.ModNone, setMainOnTop); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, 'k', gocui.ModNone, sidebarLineUp); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, 'j', gocui.ModNone, sidebarLineDown); err != nil {
		return err
	}

	return nil
}

// quit event
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// nest view event
func nextView(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		return errors.New("g or v is nil")
	}

	other := ""
	if v.Name() == mainViewName {
		other = sideViewName
	} else {
		other = mainViewName
	}

	_, err := setCurrentViewOnTop(g, other)
	return err
}

// select main view event
func setMainOnTop(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		return errors.New("g or v is nil")
	}

	_, err := setCurrentViewOnTop(g, mainViewName)
	return err
}

// select side view event
func setSideOnTop(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		return errors.New("g or v is nil")
	}

	_, err := setCurrentViewOnTop(g, sideViewName)
	return err
}

func sidebarLineUp(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		return errors.New("g or v is nil")
	}

	v.MoveCursor(0, -1, false)
	return nil
}

func sidebarLineDown(g *gocui.Gui, v *gocui.View) error {
	if g == nil || v == nil {
		return errors.New("g or v is nil")
	}
	v.MoveCursor(0, 1, false)
	return nil
}
