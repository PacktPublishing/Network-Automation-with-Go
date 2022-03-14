package main

import (
	"container/heap"
	"testing"
)

type testFlow struct {
	startCount   int
	timesSeen    int
	wantPosition int
	wantCount    int
}
type testCase struct {
	name  string
	flows map[string]testFlow
}

func TestHeap(t *testing.T) {

	for _, test := range testCases {

		h := make(Heap, 0)
		// pushing flow on the heap
		for key, f := range test.flows {
			flow := &MyFlow{
				Count: f.startCount,
				Key:   key,
			}
			heap.Push(&h, flow)

			// updating packet counts
			for j := 0; j < f.timesSeen; j++ {
				h.update(flow)
			}
		}

		// checking the results
		for i := 0; h.Len() > 0; i++ {
			f := heap.Pop(&h).(*MyFlow)

			tf := test.flows[f.Key]
			if tf.wantPosition != i {
				t.Errorf(
					"%s: unexpected position for packet key %s: got %d, want %d",
					test.name, f.Key, i, tf.wantPosition,
				)
			}

			if tf.wantCount != f.Count {
				t.Errorf(
					"%s: unexpected count for packet key %s: got %d, want %d",
					test.name, f.Key, f.Count, tf.wantCount,
				)
			}
		}
	}
}

var testCases = []testCase{
	{
		name: "single packet",
		flows: map[string]testFlow{
			"1-1": {
				startCount:   1,
				timesSeen:    0,
				wantPosition: 0,
				wantCount:    1,
			},
		},
	},
	{
		name: "last packet wins",
		flows: map[string]testFlow{
			"2-1": {
				startCount:   1,
				timesSeen:    1,
				wantPosition: 1,
				wantCount:    2,
			},
			"2-2": {
				startCount:   2,
				timesSeen:    1,
				wantPosition: 0,
				wantCount:    3,
			},
		},
	},
	{
		name: "first packet wins",
		flows: map[string]testFlow{
			"3-1": {
				startCount:   5,
				timesSeen:    0,
				wantPosition: 0,
				wantCount:    5,
			},
			"3-2": {
				startCount:   2,
				timesSeen:    0,
				wantPosition: 1,
				wantCount:    2,
			},
		},
	},
	{
		name: "last packet wins after update",
		flows: map[string]testFlow{
			"4-1": {
				startCount:   3,
				timesSeen:    0,
				wantPosition: 1,
				wantCount:    3,
			},
			"4-2": {
				startCount:   2,
				timesSeen:    2,
				wantPosition: 0,
				wantCount:    4,
			},
		},
	},
	{
		name: "first packet wins after update",
		flows: map[string]testFlow{
			"5-1": {
				startCount:   1,
				timesSeen:    3,
				wantPosition: 0,
				wantCount:    4,
			},
			"5-2": {
				startCount:   2,
				timesSeen:    0,
				wantPosition: 1,
				wantCount:    2,
			},
		},
	},
	{
		name: "tie use case/first packet wins",
		flows: map[string]testFlow{
			"6-1": {
				startCount:   2,
				timesSeen:    2,
				wantPosition: 0,
				wantCount:    4,
			},
			"6-2": {
				startCount:   3,
				timesSeen:    1,
				wantPosition: 1,
				wantCount:    4,
			},
		},
	},
	{
		name: "odd number of packets",
		flows: map[string]testFlow{
			"7-1": {
				startCount:   1,
				timesSeen:    1,
				wantPosition: 2,
				wantCount:    2,
			},
			"7-2": {
				startCount:   2,
				timesSeen:    2,
				wantPosition: 0,
				wantCount:    4,
			},
			"7-3": {
				startCount:   3,
				timesSeen:    0,
				wantPosition: 1,
				wantCount:    3,
			},
		},
	},
}

// re-using Src/DstPort field as the expected position and count values in the sorted heap
// c is current count, ep is expected position in the heap, ec is the expected count (after updates)
func testPacket(k string, c, ep, ec int) *MyFlow {
	return &MyFlow{Key: k, Count: c, SrcPort: ep, DstPort: ec}
}
