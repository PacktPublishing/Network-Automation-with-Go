package main

import (
	"container/heap"
	"testing"
)

type testCase struct {
	name    string
	packets []*MyPacket
	updates []int
}

func TestHeap(t *testing.T) {

	for _, test := range testData {

		h := make(Heap, 0)
		// pushing packets on the heap
		for _, p := range test.packets {
			heap.Push(&h, p)
		}

		// updating packet counts
		for i, update := range test.updates {
			if i > len(test.packets) {
				t.Errorf("%s: updates slice too long", test.name)
			}
			p := test.packets[i]
			for j := 0; j < update; j++ {
				h.update(p)
			}
		}

		// checking the results
		for i := 0; h.Len() > 0; i++ {
			p := heap.Pop(&h).(*MyPacket)

			// assuming position (i) == p.SrcPort
			if i != p.SrcPort {
				t.Errorf("%s: unexpected position for packet key %s: %d, expected %d", test.name, p.Key, i, p.SrcPort)
			}

			// assuming count == p.DstPort
			if p.Count != p.DstPort {
				t.Errorf("%s: unexpected count for packet key %s: %d, expected %d", test.name, p.Key, p.Count, p.DstPort)
			}
		}
	}
}

var testData = []testCase{
	{
		name:    "single packet",
		packets: []*MyPacket{testPacket("1-1", 1, 0, 1)},
		updates: []int{},
	},
	{
		name:    "last packet wins",
		packets: []*MyPacket{testPacket("2-1", 1, 1, 2), testPacket("2-2", 2, 0, 3)},
		updates: []int{1, 1},
	},
	{
		name:    "first packet wins",
		packets: []*MyPacket{testPacket("3-1", 5, 0, 5), testPacket("3-2", 2, 1, 2)},
		updates: []int{},
	},
	{
		name:    "last packet wins after update",
		packets: []*MyPacket{testPacket("4-1", 3, 1, 3), testPacket("4-2", 2, 0, 4)},
		updates: []int{0, 2},
	},
	{
		name:    "first packet wins after update",
		packets: []*MyPacket{testPacket("5-1", 1, 0, 4), testPacket("5-2", 2, 1, 2)},
		updates: []int{3, 0},
	},
	{
		name:    "tie use case/first packet wins",
		packets: []*MyPacket{testPacket("6-1", 2, 0, 4), testPacket("6-2", 3, 1, 4)},
		updates: []int{2, 1},
	},
	{
		name:    "odd number of packets",
		packets: []*MyPacket{testPacket("7-1", 1, 2, 2), testPacket("7-2", 2, 0, 4), testPacket("7-3", 3, 1, 3)},
		updates: []int{1, 2, 0},
	},
}

// re-using Src/DstPort field as the expected position and count values in the sorted heap
// c is current count, s is expected position in the heap, d is the expected count (after updates)
func testPacket(k string, c, s, d int) *MyPacket {
	return &MyPacket{Key: k, Count: c, SrcPort: s, DstPort: d}
}
