package model

import (
	"bufio"
	"io"
	"net/http"
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

// HTTPStreamFactory implements tcpassembly.StreamFactory
type HTTPStreamFactory struct {
	messages    []*Message
	messageChan chan *Message
}

// NewHTTPStreamFactory create HTTPStreamFactory object
func NewHTTPStreamFactory() *HTTPStreamFactory {
	factory := &HTTPStreamFactory{
		messages:    make([]*Message, 0, maxMessageLen),
		messageChan: make(chan *Message, maxMessageLen),
	}
	go factory.consumeMessage()
	return factory
}

func (factory *HTTPStreamFactory) consumeMessage() {
	for m := range factory.messageChan {
		factory.messages = append(factory.messages, m)
	}
	if len(factory.messages) > maxMessageLen {
		factory.messages = factory.messages[len(factory.messages)-maxMessageLen:]
	}
}

func (factory *HTTPStreamFactory) putRequest(req *http.Request) *Message {
	num := atomic.AddUint32(&seq, 1)
	m := &Message{
		Num: num,
		Req: req,
	}
	factory.messageChan <- m
	return m
}

// New create a stream object
func (factory *HTTPStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
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

// httpStream will handle the actual decoding of http requests.
type httpStream struct {
	net       gopacket.Flow
	transport gopacket.Flow
	r         tcpreader.ReaderStream
	factory   *HTTPStreamFactory
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			continue
		} else {
			tcpreader.DiscardBytesToEOF(req.Body)
			h.factory.putRequest(req)
			log.Debug("Method: %s, URL: %s, Proto: %s, ContentLength: %d, Host: %s\n", req.Method, req.URL.String(), req.Proto, req.ContentLength, req.Host)
		}
	}
}
