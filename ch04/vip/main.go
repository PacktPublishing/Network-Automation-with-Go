package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jsimonetti/rtnetlink/rtnl"
	"github.com/mdlayher/arp"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/packet"
)

const VIP1 = "198.51.100.1/32"

type vip struct {
	IP      string
	netlink *rtnl.Conn
	intf    *net.Interface
	l2Sock  *packet.Conn
}

func setupSigHandlers(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		log.Printf("Received syscall: %+v", sig)
		cancel()
	}()

}

func (c *vip) removeVIP() error {
	err := c.netlink.AddrDel(c.intf, rtnl.MustParseAddr(c.IP))
	if err != nil {
		return fmt.Errorf("could not del address: %s", err)
	}
	return nil
}

func (c *vip) addVIP() error {
	err := c.netlink.AddrAdd(c.intf, rtnl.MustParseAddr(c.IP))
	if err != nil {
		return fmt.Errorf("could not add address: %s", err)
	}
	return nil
}

func (c *vip) emitFrame(frame *ethernet.Frame) error {
	b, err := frame.MarshalBinary()
	if err != nil {
		return fmt.Errorf("error serializing frame: %s", err)
	}

	addr := &packet.Addr{HardwareAddr: ethernet.Broadcast}
	if _, err := c.l2Sock.WriteTo(b, addr); err != nil {
		return fmt.Errorf("emitFrame failed: %s", err)
	}

	log.Println("GARP sent")
	return nil
}

func (c *vip) sendGARP() error {
	ip, _, err := net.ParseCIDR(c.IP)
	if err != nil {
		return fmt.Errorf("error parsing IP: %s", err)
	}

	arpPayload, err := arp.NewPacket(
		arp.OperationReply,  // op
		c.intf.HardwareAddr, // srcHW
		ip,                  // srcIP
		c.intf.HardwareAddr, //dstHW
		ip,                  //dstIP
	)
	if err != nil {
		return fmt.Errorf("error building ARP packet: %s", err)
	}

	arpBinary, err := arpPayload.MarshalBinary()
	if err != nil {
		return fmt.Errorf("error serializing ARP packet: %s", err)
	}

	ethFrame := &ethernet.Frame{
		Destination: ethernet.Broadcast,
		Source:      c.intf.HardwareAddr,
		EtherType:   ethernet.EtherTypeARP,
		Payload:     arpBinary,
	}

	return c.emitFrame(ethFrame)
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

	ethSocket, err := packet.Listen(netIntf, packet.Raw, 0, nil)
	if err != nil {
		log.Printf("failed to ListenPacket: %v", err)
	}
	defer ethSocket.Close()

	ctx, cancel := context.WithCancel(context.Background())
	setupSigHandlers(cancel)

	v := &vip{
		IP:      VIP1,
		intf:    netIntf,
		netlink: rtnl,
		l2Sock:  ethSocket,
	}

	err = v.addVIP()
	if err != nil {
		log.Fatalf("failed to add VIP: %s", err)
	}

	timer := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			if err := v.removeVIP(); err != nil {
				log.Fatalf("failed to remove VIP: %s", err)
			}

			log.Printf("Cleanup complete")
			return
		case <-timer.C:
			if err := v.sendGARP(); err != nil {
				log.Printf("failed to send GARP: %s", err)
				cancel()
			}
		}
	}

}
