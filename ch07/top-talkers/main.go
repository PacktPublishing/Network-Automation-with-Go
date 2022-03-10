package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/netsampler/goflow2/format"
	_ "github.com/netsampler/goflow2/format/json"

	// flowpb "github.com/netsampler/goflow2/pb"
	"github.com/olekukonko/tablewriter"

	"github.com/netsampler/goflow2/utils"
	log "github.com/sirupsen/logrus"
)

const listenAddress = "sflow://:6343"

var flowMapKey = `%s:%d<->%s:%d`

type MyFlow struct {
	Key         string
	SrcAddr     string `json:"SrcAddr,omitempty"`
	DstAddr     string `json:"DstAddr,omitempty"`
	SequenceNum uint32 `json:"SequenceNum,omitempty"`
	SrcPort     int    `json:"SrcPort,omitempty"`
	DstPort     int    `json:"DstPort,omitempty"`
	ProtoName   string `json:"ProtoName,omitempty"`
	Bytes       int    `json:"Bytes,omitempty"`
	Count       int    // how many times we've received flow sample
	index       int    // The index of the item in the heap (required for update)
}

type topTalker struct {
	flowMap map[string]*MyFlow
	heap    Heap
}

type Heap []*MyFlow

func (h Heap) Len() int { return len(h) }
func (h Heap) Less(i, j int) bool {
	return h[i].Count > h[j].Count
}
func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *Heap) update(p *MyFlow) {
	p.Count++
	heap.Fix(h, p.index)
}

func (h *Heap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	n := len(*h)
	item := x.(*MyFlow)
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

	var myFlow MyFlow
	json.Unmarshal(data, &myFlow)

	if myFlow.ProtoName == "" {
		return nil
	}

	ips := []string{myFlow.SrcAddr, myFlow.DstAddr}
	sort.Strings(ips)

	var mapKey string
	if ips[0] != myFlow.SrcAddr {
		mapKey = fmt.Sprintf(
			flowMapKey,
			myFlow.SrcAddr,
			myFlow.SrcPort,
			myFlow.DstAddr,
			myFlow.DstPort,
		)
	} else {
		mapKey = fmt.Sprintf(
			flowMapKey,
			myFlow.DstAddr,
			myFlow.DstPort,
			myFlow.SrcAddr,
			myFlow.SrcPort,
		)
	}

	myFlow.Key = mapKey
	foundFlow, ok := c.flowMap[mapKey]
	if !ok {
		myFlow.Count = 1
		c.flowMap[mapKey] = &myFlow
		heap.Push(&c.heap, &myFlow)
		return nil
	}

	c.heap.update(foundFlow)

	return nil
}

func main() {

	log.Print("Top Talker app")

	ctx := context.Background()

	tt := topTalker{
		flowMap: make(map[string]*MyFlow),
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "From", "To", "Proto", "Count"})
	table.SetColMinWidth(0, 3)
	table.SetColMinWidth(1, 21)
	table.SetColMinWidth(2, 21)
	table.SetColMinWidth(3, 7)
	table.SetColMinWidth(4, 7)

	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Println("Top Talkers")
		table.ClearRows()
		for i, flow := range tt.heap {
			parts := strings.Split(flow.Key, "<->")
			table.Append([]string{fmt.Sprintf("%d", i+1), parts[0], parts[1], flow.ProtoName, fmt.Sprintf("%d", flow.Count)})
		}
		table.Render()
		time.Sleep(time.Second * 5)
	}
	fmt.Println("app exiting...")

}
