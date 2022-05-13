package main

import (
	"fmt"

	"grpc/pkg/oc"

	"github.com/openconfig/ygot/ygot"
)

func (m *Model) buildNetworkInstance(dev *oc.Device) error {
	name := defaultNetInst
	peergroup := "EBGP"
	nis := &oc.OpenconfigNetworkInstance_NetworkInstances{}
	ni, err := nis.NewNetworkInstance(name)
	if err != nil {
		return fmt.Errorf("cannot create new network instance: %w", err)
	}
	ni.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Config{
		Name: &name,
	}

	ni.Protocols = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols{}
	bgp, err := ni.Protocols.NewProtocol(bgpID, name)
	if err != nil {
		return fmt.Errorf("cannot create new bgp instance: %w", err)
	}
	bgp.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Config{
		Name:       &name,
		Identifier: bgpID,
	}

	bgp.Bgp = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp{
		Global: &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global{
			Config: &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global_Config{
				As:       ygot.Uint32(uint32(m.ASN)),
				RouterId: ygot.String(m.Loopback.IP),
			},
			AfiSafis: &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global_AfiSafis{},
		},
		PeerGroups: &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_PeerGroups{},
	}

	// Initialize the IPv4 Unicast address family.
	safi, err := bgp.Bgp.Global.AfiSafis.NewAfiSafi(ipv4uniAF)
	if err != nil {
		return fmt.Errorf("cannot enable bgp IPv4 address family: %w", err)
	}
	safi.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Global_AfiSafis_AfiSafi_Config{
		AfiSafiName: ipv4uniAF,
		Enabled:     ygot.Bool(true),
	}

	// Create Peer Group
	pg, err := bgp.Bgp.PeerGroups.NewPeerGroup(peergroup)
	if err != nil {
		return fmt.Errorf("cannot create BGP peer-group: %w", err)
	}
	pg.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_PeerGroups_PeerGroup_Config{
		PeerGroupName: ygot.String(peergroup),
	}

	pg.AfiSafis = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_PeerGroups_PeerGroup_AfiSafis{}
	pgsafi, err := pg.AfiSafis.NewAfiSafi(ipv4uniAF)
	if err != nil {
		return fmt.Errorf("cannot create BGP peer-group SAFI: %w", err)
	}
	pgsafi.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_PeerGroups_PeerGroup_AfiSafis_AfiSafi_Config{
		AfiSafiName: ipv4uniAF,
		Enabled:     ygot.Bool(true),
	}

	bgp.Bgp.Neighbors = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors{}
	for _, peer := range m.Peers {
		n, err := bgp.Bgp.Neighbors.NewNeighbor(peer.IP)
		if err != nil {
			return fmt.Errorf("cannot add bgp neighbor %s: %w", peer.IP, err)
		}
		n.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_Config{
			PeerAs:          ygot.Uint32(uint32(peer.ASN)),
			NeighborAddress: ygot.String(peer.IP),
			Enabled:         ygot.Bool(true),
			PeerGroup:       ygot.String(peergroup),
		}
	}

	if err := ni.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	dev.NetworkInstances = nis

	return nil
}
