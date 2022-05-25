package cvx

import "network.automation:input"

interface: _interfaces
router: bgp: {
	_global_bgp
}
vrf: _vrf

_global_bgp: {
	"autonomous-system": input.asn
	enable:              "on"
	"router-id":         input.loopback.ip
}

_interfaces: {
	lo: {
		ip: address: "\(input.LoopbackIP)": {}
		type: "loopback"
	}
	for intf in input.uplinks {
		"\(intf.name)": {
			type: "swp"
			ip: address: "\(intf.prefix)": {}
		}
	}
}

_vrf: {
	for vrf in input.VRFs {
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
		redistribute: connected: enable: "on"
	}
	enable: "on"
}

_neighbor: {
	for intf in input.peers {
		"\(intf.ip)": {
			type:        "numbered"
			"remote-as": intf.asn
		}
	}
}
