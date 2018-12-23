package view

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

func (ctrl *Controller) bindKey() error {
	g := ctrl.g

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

	if err := ctrl.bindMainViewKey(); err != nil {
		return err
	}

	if err := ctrl.bindSideViewKey(); err != nil {
		return err
	}

	return nil
}

func (ctrl *Controller) bindMainViewKey() error {
	g := ctrl.g
	if err := g.SetKeybinding(mainViewName, gocui.KeyArrowLeft, gocui.ModNone, setSideOnTop); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainViewName, gocui.KeyEsc, gocui.ModNone, mainEsc); err != nil {
		return err
	}
	return nil
}

func (ctrl *Controller) bindSideViewKey() error {
	g := ctrl.g
	if err := g.SetKeybinding(sideViewName, gocui.KeyArrowRight, gocui.ModNone, setMainOnTop); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, 'k', gocui.ModNone, sidebarMove(-1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, 'j', gocui.ModNone, sidebarMove(1)); err != nil {
		return err
	}

	maxY := 40
	if err := g.SetKeybinding(sideViewName, gocui.KeyCtrlF, gocui.ModNone, sidebarMove(maxY)); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, gocui.KeyCtrlB, gocui.ModNone, sidebarMove(-maxY)); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, gocui.KeyCtrlU, gocui.ModNone, sidebarMove(-maxY/2)); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, gocui.KeyCtrlD, gocui.ModNone, sidebarMove(maxY/2)); err != nil {
		return err
	}

	if err := g.SetKeybinding(sideViewName, gocui.KeyEnter, gocui.ModNone, ctrl.sidebarEnter); err != nil {
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

func sidebarMove(line int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if g == nil || v == nil {
			return errors.New("g or v is nil")
		}
		v.MoveCursor(0, line, false)
		return nil
	}
}

func (ctrl *Controller) sidebarEnter(g *gocui.Gui, v *gocui.View) error {
	log.Info("type enter")
	if g == nil || v == nil {
		return errors.New("g or v is nil")
	}
	_, y := v.Cursor()
	line, err := v.Line(y)
	log.Debugf("cursor, y: %d, line: %s", y, line)
	if err != nil {
		log.Errorf("get line failed, y: %d, err: %v", y, err)
		return nil
	}

	num := sidebarStringGetMessageNum(line)
	if num == 0 {
		log.Errorf("message num: %d", num)
		return nil
	}

	find := ctrl.mq.SearchByNum(num)
	if find == nil || find.Req == nil {
		log.Errorf("search by num, find nil")
		return nil
	}

	mainView, err := g.View(mainViewName)
	if err != nil {
		log.Errorf("find main view failed, err: %v", err)
		return nil
	}

	mainView.Clear()
	fmt.Fprintln(mainView, find.Req.Host)
	// fmt.Fprintf(mainView, "%+v\n", find.Req.Header)

	fmt.Fprintln(mainView, string(find.ReqBody))
	if find.Rsp != nil {
		fmt.Fprintln(mainView, string(find.RspBody))
	}

	return nil
}

func mainEsc(g *gocui.Gui, v *gocui.View) error {
	return nil
}
