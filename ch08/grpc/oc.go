package main

import (
	"fmt"

	"grpc/pkg/oc"

	"github.com/openconfig/ygot/ygot"
)

const (
	defaultNetInst = "default"
	policy         = "PERMIT-ALL"
	polStatemet    = "PASS"
	bgpID          = oc.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP
	ipv4uniAF      = oc.OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST
)

func (m *Model) buildNetworkInstance(dev *oc.Device) error {
	m.buildRoutePolicy(dev)
	nis := &oc.OpenconfigNetworkInstance_NetworkInstances{}
	ni, err := nis.NewNetworkInstance(defaultNetInst)
	if err != nil {
		return fmt.Errorf("cannot create new network instance: %w", err)
	}
	ni.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Config{
		Name: ygot.String(defaultNetInst),
	}

	ni.Protocols = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols{}
	bgp, err := ni.Protocols.NewProtocol(bgpID, defaultNetInst)
	if err != nil {
		return fmt.Errorf("cannot create new bgp instance: %w", err)
	}
	bgp.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Config{
		Name:       ygot.String(defaultNetInst),
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
			// PeerGroup:       ygot.String(peergroup),
		}
		n.AfiSafis = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_AfiSafis{}
		safi, err := n.AfiSafis.NewAfiSafi(ipv4uniAF)
		if err != nil {
			return fmt.Errorf("cannot add address family to bgp neighbor %s: %w", peer.IP, err)
		}
		safi.Config = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_AfiSafis_AfiSafi_Config{
			AfiSafiName: ipv4uniAF,
			Enabled:     ygot.Bool(true),
		}
		safi.ApplyPolicy = &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_AfiSafis_AfiSafi_ApplyPolicy{
			Config: &oc.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_AfiSafis_AfiSafi_ApplyPolicy_Config{
				ExportPolicy: []string{policy},
				ImportPolicy: []string{policy},
			},
		}

	}

	if err := ni.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	dev.NetworkInstances = nis

	return nil
}

func (m *Model) buildRoutePolicy(dev *oc.Device) error {
	pol := &oc.OpenconfigRoutingPolicy_RoutingPolicy{
		PolicyDefinitions: &oc.OpenconfigRoutingPolicy_RoutingPolicy_PolicyDefinitions{},
	}
	def, err := pol.PolicyDefinitions.NewPolicyDefinition(policy)
	if err != nil {
		return fmt.Errorf("cannot create policy definition %s: %w", policy, err)
	}
	def.Config = &oc.OpenconfigRoutingPolicy_RoutingPolicy_PolicyDefinitions_PolicyDefinition_Config{
		Name: ygot.String(policy),
	}
	def.Statements = &oc.OpenconfigRoutingPolicy_RoutingPolicy_PolicyDefinitions_PolicyDefinition_Statements{}
	stat, err := def.Statements.NewStatement(polStatemet)
	if err != nil {
		return fmt.Errorf("cannot create a policy statement %s: %w", polStatemet, err)
	}
	stat.Config = &oc.OpenconfigRoutingPolicy_RoutingPolicy_PolicyDefinitions_PolicyDefinition_Statements_Statement_Config{
		Name: ygot.String(polStatemet),
	}
	stat.Actions = &oc.OpenconfigRoutingPolicy_RoutingPolicy_PolicyDefinitions_PolicyDefinition_Statements_Statement_Actions{
		Config: &oc.OpenconfigRoutingPolicy_RoutingPolicy_PolicyDefinitions_PolicyDefinition_Statements_Statement_Actions_Config{
			PolicyResult: oc.OpenconfigRoutingPolicy_PolicyResultType_ACCEPT_ROUTE,
		},
	}

	dev.RoutingPolicy = pol

	return nil
}
