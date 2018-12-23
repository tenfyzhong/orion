package model

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("orion")

const (
	maxMessageLen = 100
)

type key [2]uint64

func getKey(net, transport gopacket.Flow) key {
	return [2]uint64{net.FastHash(), transport.FastHash()}
}

// HTTPStreamFactory implements tcpassembly.StreamFactory
type HTTPStreamFactory struct {
	ctlMessageChan chan *Message
	streamMap      map[key]*httpStream
	mutex          sync.RWMutex
}

// NewHTTPStreamFactory create HTTPStreamFactory object
func NewHTTPStreamFactory(ctlMessageChan chan *Message) *HTTPStreamFactory {
	factory := &HTTPStreamFactory{
		ctlMessageChan: ctlMessageChan,
		streamMap:      make(map[key]*httpStream),
	}
	return factory
}

// New create a stream object
func (factory *HTTPStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	log.Debugf("new stream, net: %+v, gopacket: %+v\n", net, transport)
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
		factory:   factory,
	}
	// Important... we must guarantee that data from the reader stream is read.
	go hstream.run()

	return &hstream.r
}

func (factory *HTTPStreamFactory) putMessage(m *Message) {
	num := atomic.AddUint32(&seq, 1)
	m.Num = num
	factory.ctlMessageChan <- m
}

func (factory *HTTPStreamFactory) getStream(k key) *httpStream {
	factory.mutex.RLock()
	defer factory.mutex.RUnlock()
	return factory.streamMap[k]
}

func (factory *HTTPStreamFactory) setStream(k key, h *httpStream) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()
	if h != nil {
		factory.streamMap[k] = h
	} else {
		delete(factory.streamMap, k)
	}
}

// httpStream will handle the actual decoding of http requests.
type httpStream struct {
	net       gopacket.Flow
	transport gopacket.Flow
	r         tcpreader.ReaderStream
	factory   *HTTPStreamFactory
	message   *Message
}

const (
	linkTypeUnknow = iota
	linkTypeReq
	linkTypeRsp
)

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	linkType := 0
	for {
		if linkType == 0 {
			linkType = getLinkType(buf)
			if linkType == 0 {
				continue
			}
		}

		if linkType == 1 {
			err := h.processRequest(buf)
			if err == io.EOF {
				// We must read until we see an EOF... very important!
				return
			}
		} else if linkType == 2 {
			err := h.processResponse(buf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				// We must read until we see an EOF... very important!
				return
			}
		} else {
			panic("wtf")
		}
	}
}

func getLinkType(buf *bufio.Reader) int {
	first, err := buf.Peek(5)
	if err != nil {
		return 0
	}

	log.Debug("first: ", string(first))

	if strings.HasPrefix(string(first), "HTTP/") {
		return linkTypeRsp
	}
	return linkTypeReq

}

func (h *httpStream) processRequest(buf *bufio.Reader) error {
	req, err := http.ReadRequest(buf)
	if err != nil {
		log.Info("read error", err)
		return err
	}

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	log.Debugf("Method: %s, URL: %s, Proto: %s, ContentLength: %d, Host: %s, bodylen: %d\n", req.Method, req.URL.String(), req.Proto, req.ContentLength, req.Host, len(body))

	k := getKey(h.net, h.transport)
	storeStream := h.factory.getStream(k)
	if storeStream == nil {
		m := &Message{
			Req:     req,
			ReqBody: body,
		}
		h.message = m
		h.factory.putMessage(m)
		h.factory.setStream(k, h)
	} else {
		m := storeStream.message
		m.Req = req
		m.ReqBody = body
		h.factory.putMessage(m)
		h.factory.setStream(k, nil)
	}

	return nil
}

func (h *httpStream) processResponse(buf *bufio.Reader) error {
	rsp, err := http.ReadResponse(buf, nil)
	if err != nil {
		log.Info("read error", err)
		return err
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	log.Debugf("read body, status: %d, len: %d\n", rsp.StatusCode, len(body))
	rsp.Body.Close()

	k := getKey(h.net, h.transport)
	storeStream := h.factory.getStream(k)
	if storeStream == nil {
		m := &Message{
			Rsp:     rsp,
			RspBody: body,
		}
		h.message = m
		h.factory.setStream(k, h)
	} else {
		m := storeStream.message
		m.Rsp = rsp
		m.RspBody = body
		h.factory.setStream(k, nil)
	}

	return nil
}
