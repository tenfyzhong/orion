package main

import (
	"flag"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/tenfyzhong/orion/controller"
	"github.com/tenfyzhong/orion/view"
)

var iface = flag.String("i", "", "Interface to get packets from")
var fname = flag.String("r", "", "Filename to read from, overrides -i")
var snaplen = flag.Int("s", 1600, "SnapLen for pcap packet capture")
var filter = flag.String("f", "tcp and dst port 80", "BPF filter for pcap")
var verbose = flag.Bool("v", false, "Logs every packet in great detail")

func main() {
	flag.Parse()
	if *iface == "" {
		*iface = "any"
	}
	if *filter == "" {
		*filter = "host localhost"
	}

	ctl := controller.NewController(*iface, *snaplen, *filter)
	err := ctl.Init()
	if err != nil {
		log.Critical("init failed, ", err)
		os.Exit(-1)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Critical("new gui, ", err)
		os.Exit(-2)
	}

	viewController := view.NewController(g, 1024)
	err = viewController.Run()
	if err != nil {
		log.Critical("run view, ", err)
		os.Exit(-3)
	}

	// g, err := view.Run()
	// if err != nil {
	// 	return
	// }

	ctl.AddUpdateFunc(viewController.Update)

	go ctl.Run()

	// view.InitTitle(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return
	}

	g.Close()
}
