package cvx

import (
	"net"
)

_input: _Input

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
	}]
	VRFs: [{name: "default"}]
	PeerGroup: ""
}

_nvue: {
	interface: _interfaces
	router: bgp: {
		_global_bgp
	}
	vrf: _vrf
}

_global_bgp: {
	"autonomous-system": _input.ASN
	enable:              "on"
	"router-id":         _input.RouterID
}

_interfaces: {
	lo: {
		ip: address: "\(_input.LoopbackIP)": {}
		type: "loopback"
	}
	for intf in _input.Uplinks {
		"\(intf.name)": {
			type: "swp"
			ip: address: "\(intf.ip)/\(intf.prefixLen)": {}
		}
	}
}

_vrf: {
	for vrf in _input.VRFs {
		"\(vrf.name)": {
			router: bgp: _vrf_bgp
			if vrf.name == "default" {
				router: bgp: neighbor: _neighbor
			}
		}
	}
}

_vrf_bgp: {
	"address-family": "ipv4-unicast": {
		enable: "on"
		redistribute: connected: enable: "on"
	}
	enable: "on"
}

_neighbor: {
	for intf in _input.Peers {
		"\(intf.ip)": {
			type:        string | *"numbered"
			"remote-as": "\(intf.asn)"
		}
	}
}

