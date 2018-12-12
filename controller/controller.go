package controller

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/tenfyzhong/orion/model"
)

// Controller controller model
type Controller struct {
	iface         string
	snaplen       int
	filter        string
	handle        *pcap.Handle
	streamFactory *model.HTTPStreamFactory
}

// NewController create Controller object
func NewController(iface string, snaplen int, filter string) *Controller {
	return &Controller{
		iface:         iface,
		snaplen:       snaplen,
		filter:        filter,
		streamFactory: model.NewHTTPStreamFactory(),
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
		return err
	}

	if err := handle.SetBPFFilter(c.filter); err != nil {
		return err
	}

	c.handle = handle
	return nil
}

// Run begin capture packet
func (c *Controller) Run() {
	streamPool := tcpassembly.NewStreamPool(c.streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("reading in packets")

	// read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates end of a pcap file.
			if packet == nil {
				return
			}

			if packet.NetworkLayer() == nil ||
				packet.TransportLayer() == nil ||
				packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
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
