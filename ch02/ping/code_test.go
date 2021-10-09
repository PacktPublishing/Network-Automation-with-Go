package ping_test

import (
	"testing"
	"github.com/PacktPublishing/Network-Automation-with-Go/ch02/ping"
)

func TestSend(t *testing.T) {
	cases := []struct {
		want   string
	}{
		{want: "pong"},
	}
	for _, c := range cases {
		result := ping.Send()
		if result != c.want {
			t.Fatalf("[%s] is incorrect, we want [%s]", result, c.want)
		}
	}
}