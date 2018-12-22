package model

import "net/http"

var seq uint32

// Message a package message
type Message struct {
	Num     uint32
	Req     *http.Request
	Rsp     *http.Response
	ReqBody []byte
}
