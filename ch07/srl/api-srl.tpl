# https://go.dev/play/p/TwG8sUkmLcs
- path: "/interface[name=system0]/subinterface[index=0]/ipv4/address[ip-prefix={{ .Loopback.IP }}/32]"
  value: '{}'

- path: "/network-instance[name=default]/interface[name=system0.0]"
  value: '{}'

- path: "/interface[name={{ (index .Uplinks 0).Name }}]/subinterface[index=0]/ipv4/address[ip-prefix={{ (index .Uplinks 0).Prefix }}]"
  value: '{}'

- path: "/network-instance[name=default]/interface[name={{ (index .Uplinks 0).Name }}.0]"
  value: '{}'

- path: "/interface[name={{ (index .Uplinks 0).Name }}]/admin-state"
  value: enable

- path: "/routing-policy/policy[name=all]/default-action/accept/bgp/local-preference/set"
  value: '100'

- path: "/network-instance[name=default]/protocols/bgp"
  value: '{ "autonomous-system": "{{ .ASN }}","router-id": "{{ .Loopback.IP }}" }'

- path: "/network-instance[name=default]/protocols/bgp/ipv4-unicast/admin-state"
  value: enable

- path: "/network-instance[name=default]/protocols/bgp/group[group-name=EBGP]"
  value: '{"export-policy": "all","import-policy": "all"}'

- path: "/network-instance[name=default]/protocols/bgp/neighbor[peer-address={{ (index .Peers 0).IP }}]"
  value: '{"peer-as":{{ (index .Peers 0).ASN }},"peer-group": "EBGP"}'
