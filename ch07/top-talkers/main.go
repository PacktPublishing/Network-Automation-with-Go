package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/netsampler/goflow2/format"
	_ "github.com/netsampler/goflow2/format/json"

	// flowpb "github.com/netsampler/goflow2/pb"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/netsampler/goflow2/utils"
	log "github.com/sirupsen/logrus"
)

const listenAddress = "sflow://:6343"

var flowMapKey = `%s:%d<->%s:%d`

type MyPacket struct {
	Key         string
	SrcAddr     string `json:"SrcAddr,omitempty"`
	DstAddr     string `json:"DstAddr,omitempty"`
	SequenceNum uint32 `json:"SequenceNum,omitempty"`
	SrcPort     int    `json:"SrcPort,omitempty"`
	DstPort     int    `json:"DstPort,omitempty"`
	ProtoName   string `json:"ProtoName,omitempty"`
	Count       int    // how many times we;ve received flow sample
	index       int    // The index of the item in the heap (required for update)
}

type topTalker struct {
	flowMap map[string]*MyPacket
	heap    Heap
}

type Heap []*MyPacket

func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].Count > h[j].Count }
func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *Heap) update(p *MyPacket) {
	p.Count++
	heap.Fix(h, p.index)
}

func (h *Heap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	n := len(*h)
	item := x.(*MyPacket)
	item.index = n
	*h = append(*h, item)

}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // avoid memory leak
	x.index = -1   // for safety
	*h = old[0 : n-1]
	return x
}

func (c *topTalker) Send(key, data []byte) error {
	//log.Printf("transport.Send : Key %s, data %+v\n", key, string(data))

	var myPacket MyPacket
	json.Unmarshal(data, &myPacket)

	ips := []string{myPacket.SrcAddr, myPacket.DstAddr}
	sort.Strings(ips)

	var mapKey string
	if ips[0] != myPacket.SrcAddr {
		mapKey = fmt.Sprintf(flowMapKey, myPacket.SrcAddr, myPacket.SrcPort, myPacket.DstAddr, myPacket.DstPort)
	} else {
		mapKey = fmt.Sprintf(flowMapKey, myPacket.DstAddr, myPacket.DstPort, myPacket.SrcAddr, myPacket.SrcPort)
	}

	myPacket.Key = mapKey
	i, ok := c.flowMap[mapKey]
	if !ok {
		myPacket.Count = 1
		c.flowMap[mapKey] = &myPacket
		heap.Push(&c.heap, &myPacket)
		return nil
	}

	c.heap.update(i)

	return nil
}

func main() {

	log.Print("Top Talker app")

	ctx := context.Background()

	tt := topTalker{
		flowMap: make(map[string]*MyPacket),
		heap:    make(Heap, 0),
	}

	heap.Init(&tt.heap)

	formatter, err := format.FindFormat(ctx, "json")
	if err != nil {
		log.Fatal(err)
	}

	listenAddrUrl, err := url.Parse(listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	hostname := listenAddrUrl.Hostname()
	port, err := strconv.ParseUint(listenAddrUrl.Port(), 10, 64)
	if err != nil {
		log.Errorf("Port %s could not be converted to integer", listenAddrUrl.Port())
		return
	}

	sSFlow := &utils.StateSFlow{
		Format:    formatter,
		Logger:    log.StandardLogger(),
		Transport: &tt,
	}

	go sSFlow.FlowRoutine(1, hostname, int(port), false)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	table := widgets.NewTable()
	table.BorderStyle = ui.NewStyle(ui.ColorGreen)
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 120, 10)

	go func() {
		for {

			table.Rows = [][]string{
				[]string{"Position", "From", "To", "Proto", "Count"},
			}
			for i, flow := range tt.heap {
				table.Rows = append(table.Rows, []string{
					fmt.Sprintf("%d", i+1),
					fmt.Sprintf("%s:%d", flow.SrcAddr, flow.SrcPort),
					fmt.Sprintf("%s:%d", flow.DstAddr, flow.DstPort),
					flow.ProtoName,
					fmt.Sprintf("%d", flow.Count),
				})
			}

			ui.Render(table)
			time.Sleep(time.Second * 1)
		}
	}()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}

}
