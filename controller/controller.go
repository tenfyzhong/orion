package controller

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	logging "github.com/op/go-logging"
	"github.com/tenfyzhong/orion/model"
)

var log = logging.MustGetLogger("orion")

// UpdateFunc update callback
type UpdateFunc func(m *model.Message)

// Controller controller model
type Controller struct {
	iface         string
	snaplen       int
	filter        string
	handle        *pcap.Handle
	streamFactory *model.HTTPStreamFactory
	updateFuncs   []UpdateFunc
	messageChan   chan *model.Message
}

// NewController create Controller object
func NewController(iface string, snaplen int, filter string) *Controller {
	log.Info("iface: %s, snaplen: %d, filter: %s\n", iface, snaplen, filter)
	messageChan := make(chan *model.Message, 100)
	c := &Controller{
		iface:         iface,
		snaplen:       snaplen,
		filter:        filter,
		streamFactory: model.NewHTTPStreamFactory(messageChan),
		updateFuncs:   make([]UpdateFunc, 0, 0),
		messageChan:   messageChan,
	}
	go c.consumeMessage()
	return c
}

func (c *Controller) consumeMessage() {
	for m := range c.messageChan {
		for _, f := range c.updateFuncs {
			f(m)
		}
	}
}

// Init init the controller
func (c *Controller) Init() error {
	handle, err := pcap.OpenLive(
		c.iface,
		int32(c.snaplen),
		true,
		pcap.BlockForever)
	if err != nil {
		log.Critical("open live failed", err)
		return err
	}

	if err := handle.SetBPFFilter(c.filter); err != nil {
		log.Critical("set bpf filter failed", err)
		return err
	}

	c.handle = handle

	return nil
}

// AddUpdateFunc Add updateFunc
func (c *Controller) AddUpdateFunc(u UpdateFunc) {
	c.updateFuncs = append(c.updateFuncs, u)
}

// Run begin capture packet
func (c *Controller) Run() {
	streamPool := tcpassembly.NewStreamPool(c.streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Info("reading in packets")

	// read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates end of a pcap file.
			if packet == nil {
				log.Info("get a nil packet")
				return
			}

			if packet.NetworkLayer() == nil ||
				packet.TransportLayer() == nil ||
				packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Warning("Unusable packet")
				continue
			}

			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(
				packet.NetworkLayer().NetworkFlow(),
				tcp,
				packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seend activity in
			// the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
