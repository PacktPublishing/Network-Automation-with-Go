package main

import (
	"fmt"

	"grpc/pkg/oc"

	"github.com/openconfig/ygot/ygot"
)

const (
	defaultNetInst     = "default"
	policy             = "PERMIT-ALL"
	polStatemet        = "PASS"
	subscriptionIDName = "BGP"
	sensorGroupID      = "BGPNeighbor"
	subsInterval       = 2000
	bgpID              = oc.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP
	ipv4uniAF          = oc.OpenconfigBgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST
	subsPath           = "Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor"
)

func (m *Model) buildNetworkInstance(dev *oc.Device) error {
	m.buildTelemetrySubs(dev)
	m.buildRoutePolicy(dev)
	nis := new(oc.OpenconfigNetworkInstance_NetworkInstances)
	ni, err := nis.NewNetworkInstance(defaultNetInst)
	if err != nil {
		return fmt.Errorf(
			"cannot create new network instance: %w",
			err,
		)
	}
	ygot.BuildEmptyTree(ni)
	ni.Config.Name = ygot.String(defaultNetInst)

	ygot.BuildEmptyTree(ni.Protocols)
	bgp, err := ni.Protocols.NewProtocol(bgpID, defaultNetInst)
	if err != nil {
		return fmt.Errorf("cannot create new bgp instance: %w", err)
	}
	ygot.BuildEmptyTree(bgp)

	bgp.Config.Name = ygot.String(defaultNetInst)
	bgp.Config.Identifier = bgpID
	bgp.Bgp.Global.Config.As = ygot.Uint32(uint32(m.ASN))
	bgp.Bgp.Global.Config.RouterId = ygot.String(m.Loopback.IP)

	safi, err := bgp.Bgp.Global.AfiSafis.NewAfiSafi(ipv4uniAF)
	if err != nil {
		return fmt.Errorf(
			"cannot enable bgp IPv4 address family: %w",
			err,
		)
	}
	ygot.BuildEmptyTree(safi)

	safi.Config.AfiSafiName = ipv4uniAF
	safi.Config.Enabled = ygot.Bool(true)

	for _, peer := range m.Peers {
		n, err := bgp.Bgp.Neighbors.NewNeighbor(peer.IP)
		if err != nil {
			return fmt.Errorf(
				"cannot add bgp neighbor %s: %w",
				peer.IP,
				err,
			)
		}
		ygot.BuildEmptyTree(n)

		n.Config.PeerAs = ygot.Uint32(uint32(peer.ASN))
		n.Config.NeighborAddress = ygot.String(peer.IP)
		n.Config.Enabled = ygot.Bool(true)
		// n.Config.PeerGroup = ygot.String(peergroup)

		safi, err := n.AfiSafis.NewAfiSafi(ipv4uniAF)
		if err != nil {
			return fmt.Errorf(
				"cannot add address family to bgp neighbor %s: %w",
				peer.IP,
				err,
			)
		}
		ygot.BuildEmptyTree(safi)
		safi.Config.AfiSafiName = ipv4uniAF
		safi.Config.Enabled = ygot.Bool(true)

		safi.ApplyPolicy.Config.ExportPolicy = []string{policy}
		safi.ApplyPolicy.Config.ImportPolicy = []string{policy}

	}

	if err := ni.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	dev.NetworkInstances = nis

	return nil
}

func (m *Model) buildRoutePolicy(dev *oc.Device) error {
	pol := new(oc.OpenconfigRoutingPolicy_RoutingPolicy)
	ygot.BuildEmptyTree(pol)

	def, err := pol.PolicyDefinitions.NewPolicyDefinition(policy)
	if err != nil {
		return fmt.Errorf(
			"cannot create policy definition %s: %w",
			policy,
			err,
		)
	}
	ygot.BuildEmptyTree(def)
	def.Config.Name = ygot.String(policy)

	stat, err := def.Statements.NewStatement(polStatemet)
	if err != nil {
		return fmt.Errorf(
			"cannot create a policy statement %s: %w",
			polStatemet,
			err,
		)
	}
	ygot.BuildEmptyTree(stat)
	stat.Config.Name = ygot.String(polStatemet)
	stat.Actions.Config.PolicyResult = oc.OpenconfigRoutingPolicy_PolicyResultType_ACCEPT_ROUTE

	dev.RoutingPolicy = pol

	return nil
}

/*
    "error-path": "openconfig-telemetry:telemetry-system/",
    "error-message": "Unknown element is specified.",
    "error-info": {
    "bad-element": "persistent-subscriptions"

FIX: Modify OC YANG model :-(
*/
func (m *Model) buildTelemetrySubs(dev *oc.Device) error {
	t := new(oc.OpenconfigTelemetry_TelemetrySystem)
	ygot.BuildEmptyTree(t)

	sg, err := t.SensorGroups.NewSensorGroup(sensorGroupID)
	if err != nil {
		return fmt.Errorf(
			"failed to generate sensor group %s: %w",
			sensorGroupID,
			err,
		)
	}
	ygot.BuildEmptyTree(sg)
	sg.Config.SensorGroupId = ygot.String(sensorGroupID)

	ygot.BuildEmptyTree(sg)
	sp, err := sg.SensorPaths.NewSensorPath(subsPath)
	if err != nil {
		return fmt.Errorf(
			"failed to generate sensor path %s: %w",
			subsPath,
			err,
		)
	}
	ygot.BuildEmptyTree(sp)
	sp.Path = ygot.String(subsPath)
	sp.Config.Path = ygot.String(subsPath)

	// Without Modified OpenConfig model
	sb, err := t.Subscriptions.Persistent.NewSubscription(
		subscriptionIDName,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to generate telemetry subscription %v: %w",
			subscriptionIDName,
			err,
		)
	}
	ygot.BuildEmptyTree(sb)
	sb.Config.SubscriptionId = ygot.String(subscriptionIDName)
	sp.Path = ygot.String(subsPath)
	sp.Config.Path = ygot.String(subsPath)

	spf, err := sb.SensorProfiles.NewSensorProfile(sensorGroupID)
	if err != nil {
		return fmt.Errorf(
			"failed to generate sensor profile %s: %w",
			sensorGroupID,
			err,
		)
	}
	ygot.BuildEmptyTree(spf)
	spf.Config.SensorGroup = ygot.String(sensorGroupID)
	spf.Config.SampleInterval = ygot.Uint64(subsInterval)

	dev.TelemetrySystem = t

	return nil
}
