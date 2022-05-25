package input

import (
	"net"
)

asn: <=65535 & >=64512
loopback: ip: net.IPv4 & string
uplinks: [...{
	name:   string
	prefix: net.IPCIDR & string
}]
peers: [...{
	ip:  net.IPv4 & string
	asn: <=65535 & >=64512
}]
LoopbackIP: "\(loopback.ip)/32"
VRFs: [{name: "default"}]
