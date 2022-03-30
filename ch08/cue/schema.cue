package cvx

import (
	"net"
)

_Input: {
	ASN:        <=65535 & >=64512
	RouterID:   net.IPv4 & string
	LoopbackIP: "\(RouterID)/32"
	Uplinks: [...{
		name:      string
		ip:        net.IPv4 & string
		prefixLen: <=31 & >=1
	}]
	Peers: [...{
		ip: net.IPv4 & string
		asn: <=65535 & >=64512
	}]
	VRFs: [{name: "default"}]
	PeerGroup: ""
}