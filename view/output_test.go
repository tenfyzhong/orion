package view

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenfyzhong/orion/model"
)

func TestSidebarTitle(t *testing.T) {
	assert.Equal(t, "Num Method Host Status", sidebarTitle())
}

func messageSidebarNil(t *testing.T) {
	assert.Equal(t, "", messageSidebarString(nil))
	assert.Equal(t, "", messageSidebarString(&model.Message{
		Num: 1,
	}))
}

func TestMessageSidebarString(t *testing.T) {
	req, err := http.NewRequest("GET", "http://tenfy.cn", nil)
	assert.NoError(t, err)
	m := &model.Message{
		Num: 1,
		Req: req,
		Rsp: nil,
	}
	assert.Equal(t, "  1    GET tenfy.cn Pending", messageSidebarString(m))
}

func TestMessageSidebarStringWithRsp(t *testing.T) {
	req, err := http.NewRequest("GET", "http://tenfy.cn", nil)
	assert.NoError(t, err)
	m := &model.Message{
		Num: 1,
		Req: req,
		Rsp: &http.Response{},
	}
	assert.Equal(t, "  1    GET tenfy.cn OK", messageSidebarString(m))
}

func TestSidebarStringGetMessageNumEmpty(t *testing.T) {
	assert.Equal(t, 0, sidebarStringGetMessageNum(""))
	assert.Equal(t, 0, sidebarStringGetMessageNum(" "))
	assert.Equal(t, 0, sidebarStringGetMessageNum("G "))
	assert.Equal(t, 0, sidebarStringGetMessageNum(" G"))
}

func TestSidebarStringGetMessageNum(t *testing.T) {
	assert.Equal(t, 1, sidebarStringGetMessageNum("  1 GET"))
	assert.Equal(t, 1, sidebarStringGetMessageNum("1 GET"))
	assert.Equal(t, 1, sidebarStringGetMessageNum(" 1 1"))
}
