package view

import (
	"fmt"

	"github.com/tenfyzhong/orion/model"
)

func sidebarTitle() string {
	return "Num Method Host Status"
}

func messageSidebarString(m *model.Message) string {
	if m == nil || m.Req == nil {
		return ""
	}
	// num method host status
	if m.Rsp == nil {
		return fmt.Sprintf("%3d %6s %s %s", m.Num, m.Req.Method, m.Req.Host, "Pending")
	}
	return fmt.Sprintf("%3d %6s %s %s", m.Num, m.Req.Method, m.Req.Host, "OK")
}
