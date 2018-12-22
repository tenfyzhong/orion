package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/tenfyzhong/orion/model"
)

// MessageController controller the message to draw
type MessageController struct {
	g        *gocui.Gui
	mq       *MessageQueue
	capacity int
}

// NewMessageController new MessageController
func NewMessageController(g *gocui.Gui, capture int) *MessageController {
	return &MessageController{
		g:  g,
		mq: NewMessageQueue(capture),
	}
}

// Update update callback
func (mc *MessageController) Update(m *model.Message) {
	if m == nil {
		return
	}

	find := mc.mq.SearchByNum(m.Num)
	if find != nil {
		if m.Req != nil {
			find.Req = m.Req
		}
		if m.Rsp != nil {
			find.Rsp = m.Rsp
		}
	}
	mc.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(sideViewName)
		if err != nil {
			log.Error("get view failed", err)
			return err
		}

		if find == nil {
			appendMessage(v, m)
		} else {
			updateMessage(v, find)
		}
		return nil
	})
}

func appendMessage(v *gocui.View, m *model.Message) {
	if v == nil || m == nil {
		return
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
	// v.SetCursor(0, newY)
	v.MoveCursor(0, 1, false)
}

func updateMessage(v *gocui.View, m *model.Message) {
	if v == nil || m == nil {
		return
	}
	lines := v.BufferLines()
	index := getLineIndexStartWithNum(lines, m.Num)
	if index < 0 {
		return
	}
	line := lines[index]
	clearLine(v, index, line)
	newLine := messageSidebarString(m)
	writeLine(v, index, newLine)
}

func clearLine(v *gocui.View, index int, line string) {
	xLen := len([]byte(line))
	var empty rune
	for i := 1; i <= xLen; i++ {
		v.SetCursor(i, index)
		v.EditWrite(empty)
	}
}

func writeLine(v *gocui.View, index int, line string) {
	for i, c := range []rune(line) {
		v.SetCursor(i+1, index)
		v.EditWrite(c)
	}
}

func getLineIndexStartWithNum(lines []string, num uint32) int {
	if len(lines) == 0 {
		return -1
	}

	findIndex := 0
	for {
		if findIndex < 0 || findIndex >= len(lines) {
			return -1
		}
		str := lines[findIndex]
		lineNum := sidebarStringGetMessageNum(str)
		if lineNum < num {
			findIndex += int(num - lineNum)
		} else if lineNum > num {
			findIndex -= int(lineNum - num)
		} else {
			return findIndex
		}
	}
}
