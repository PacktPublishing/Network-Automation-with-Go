package cvx

_input: _Input & {
	ASN:      65002
	RouterID: "198.51.100.2"
	Uplinks: [{
		name:      "swp1"
		ip:        "192.0.2.3"
		prefixLen: 31
	}]
	Peers: [{
		ip:  "192.0.2.2"
		asn: 65001
	}]
}

