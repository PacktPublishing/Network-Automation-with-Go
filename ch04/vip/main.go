package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jsimonetti/rtnetlink/rtnl"
	"github.com/mdlayher/arp"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

const VIP1 = "198.51.100.1/32"

type vip struct {
	IP      string
	netlink *rtnl.Conn
	intf    *net.Interface
	raw     *raw.Conn
	cancelF context.CancelFunc
}

func newVIP(ip string, intf *net.Interface, nl *rtnl.Conn, raw *raw.Conn, cf context.CancelFunc) *vip {
	return &vip{
		IP:      ip,
		intf:    intf,
		netlink: nl,
		raw:     raw,
		cancelF: cf,
	}
}

func (c *vip) setupSigHandlers() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		log.Printf("Received syscall:%+v", sig)
		c.cancelF()
	}()

}

func (c *vip) removeVIP() {
	err := c.netlink.AddrDel(c.intf, rtnl.MustParseAddr(c.IP))
	if err != nil {
		log.Printf("could not del address: %s", err)
	}
}

func (c *vip) addVIP() {
	err := c.netlink.AddrAdd(c.intf, rtnl.MustParseAddr(c.IP))
	if err != nil {
		log.Printf("could not add address: %s", err)
		c.cancelF()
	}
}

func (c *vip) emitFrame(frame *ethernet.Frame) {
	b, err := frame.MarshalBinary()
	if err != nil {
		log.Printf("failed to marshal frame: %s", err)
		return
	}

	addr := &raw.Addr{HardwareAddr: ethernet.Broadcast}
	if _, err := c.raw.WriteTo(b, addr); err != nil {
		log.Printf("emitFrame failed: %s", err)
	}
	log.Printf("GARP sent: %+v", frame)
}

func (c *vip) sendGARP() {
	ip, _, err := net.ParseCIDR(c.IP)
	if err != nil {
		log.Printf("error parsing IP: %s", err)
	}

	arpPayload, err := arp.NewPacket(
		arp.OperationReply,  // op
		c.intf.HardwareAddr, // srcHW
		ip,                  // srcIP
		c.intf.HardwareAddr, //dstHW
		ip,                  //dstIP
	)
	if err != nil {
		log.Printf("arpPayload: %s", err)
	}

	arpBinary, err := arpPayload.MarshalBinary()
	if err != nil {
		log.Printf("arpBinary: %s", err)
	}

	ethFrame := &ethernet.Frame{
		Destination: ethernet.Broadcast,
		Source:      c.intf.HardwareAddr,
		EtherType:   ethernet.EtherTypeARP,
		Payload:     arpBinary,
	}

	c.emitFrame(ethFrame)

}

func main() {
	intfStr := flag.String("intf", "", "VIP interface")
	flag.Parse()

	if *intfStr == "" {
		log.Fatal("Please provide -intf flag")
	}

	netIntf, err := net.InterfaceByName(*intfStr)
	if err != nil {
		log.Fatalf("interface not found: %s", err)
	}

	rtnl, err := rtnl.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer rtnl.Close()

	raw, err := raw.ListenPacket(netIntf, uint16(ethernet.EtherTypeARP), nil)
	if err != nil {
		log.Printf("failed to ListenPacket: %v", err)
	}
	defer raw.Close()

	ctx, cancel := context.WithCancel(context.Background())

	v := newVIP(VIP1, netIntf, rtnl, raw, cancel)

	v.setupSigHandlers()

	v.addVIP()

	for {
		select {
		case <-ctx.Done():
			v.removeVIP()
			return
		default:
			v.sendGARP()
			time.Sleep(3 * time.Second)
		}
	}

}
