package view

import (
	"net/http"
	"testing"

	"github.com/tenfyzhong/orion/model"
)

func TestSidebarTitle(t *testing.T) {
	if sidebarTitle() != "Num Method Host Status" {
		t.FailNow()
	}
}

func TestSidebarStringMNil(t *testing.T) {
	if messageSidebarString(nil) != "" {
		t.FailNow()
	}
}

func TestSidebarStringReqNil(t *testing.T) {
	m := &model.Message{}
	if messageSidebarString(m) != "" {
		t.FailNow()
	}
}

func TestSidebarStringRspNil(t *testing.T) {
	req, err := http.NewRequest("GET", "https://www.tenfy.cn", nil)
	if err != nil || req == nil {
		t.Fatal("err should be nil, req should not be nil", err, req)
	}
	m := &model.Message{
		Num: 1,
		Req: req,
	}
	expect := "  1    GET www.tenfy.cn Pending"
	actual := messageSidebarString(m)
	if actual != expect {
		t.Fatalf("expect: %s, but get: %s", expect, actual)
	}
}

func TestSidebarString(t *testing.T) {
	req, err := http.NewRequest("GET", "https://www.tenfy.cn", nil)
	if err != nil || req == nil {
		t.Fatal("err should be nil, req should not be nil", err, req)
	}
	m := &model.Message{
		Num: 1,
		Req: req,
		Rsp: &http.Response{},
	}
	expect := "  1    GET www.tenfy.cn OK"
	actual := messageSidebarString(m)
	if actual != expect {
		t.Fatalf("expect: %s, but get: %s", expect, actual)
	}
}
