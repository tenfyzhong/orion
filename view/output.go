package view

import (
	"fmt"
	"unicode"

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

func sidebarStringGetMessageNum(str string) uint32 {
	runes := []rune(str)
	if len(runes) == 0 {
		return 0
	}
	i := 0

	// trim spaces
	for i < len(runes) && runes[i] == ' ' {
		i++
		continue
	}

	// parse num
	result := uint32(0)
	for i < len(runes) && unicode.IsNumber(runes[i]) {
		result = result*10 + uint32(runes[i]-'0')
		i++
	}

	return result
}
