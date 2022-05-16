# gRPC example

## Generating Go binding

### From protobuf files

The Go generated code in `ems_grpc.pb.go` is the result of the following:

- `proto/ems`

```bash
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --go_opt=Mproto/ems/ems_grpc.proto=proto/ems \
    --go-grpc_opt=Mproto/ems/ems_grpc.proto=proto/ems \
    proto/ems/ems_grpc.proto
```

- `proto/telemetry`

The Go generated code in `oc.go` is the result of the following:

```bash
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --go_opt=Mproto/telemetry/telemetry.proto=proto/telemetry \
    --go-grpc_opt=Mproto/telemetry/telemetry.proto=proto/telemetry \
    proto/telemetry/telemetry.proto
```

### From YANG models

- `pkg/oc`

The Go generated code in `oc.go` is the result of the following:

```bash
generator -path=yang \
  -generate_fakeroot \
  -fakeroot_name=device \
  -output_file=pkg/oc/oc.go \
  -package_name=oc \
  -exclude_modules=ietf-interfaces \
  yang/openconfig/network-instance/openconfig-network-instance.yang \
  yang/openconfig/interfaces/openconfig-interfaces.yang \
  yang/openconfig/telemetry/openconfig-telemetry-modified.yang
```

We use a slightly modified version of the `openconfig-telemetry` model. For details see: https://github.com/openconfig/public/issues/647

## Config

```json
$ go run decode-telemetry


BGP config applied on sandbox-iosxr-1.cisco.com:57777


Config from sandbox-iosxr-1.cisco.com:
{
 "openconfig-network-instance:network-instances": {
  "network-instance": [
   {
    "name": "default",
    "protocols": {
     "protocol": [
      {
       "identifier": "openconfig-policy-types:BGP",
       "name": "default",
       "bgp": {
        "global": {
         "config": {
          "as": 65000,
          "router-id": "198.51.100.0"
         },
         "afi-safis": {
          "afi-safi": [
           {
            "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
            "config": {
             "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
             "enabled": true
            }
           }
          ]
         }
        },
        "neighbors": {
         "neighbor": [
          {
           "neighbor-address": "192.0.2.1",
           "config": {
            "neighbor-address": "192.0.2.1",
            "peer-as": 65001,
            "enabled": true
           },
           "afi-safis": {
            "afi-safi": [
             {
              "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
              "config": {
               "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
               "enabled": true
              },
              "apply-policy": {
               "config": {
                "import-policy": [
                 "PERMIT-ALL"
                ],
                "export-policy": [
                 "PERMIT-ALL"
                ]
               }
              }
             }
            ]
           }
          }
         ]
        }
       }
      }
     ]
    }
   }
  ]
 }
}


----
Time: Mon May 16 13:58:20 2022
Path: Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor

  Neighbor:  192.0.2.1
  Connection state:  bgp-st-idle

----
Time: Mon May 16 13:58:22 2022
Path: Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor

  Neighbor:  192.0.2.1
  Connection state:  bgp-st-idle

----
Time: Mon May 16 13:58:24 2022
Path: Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor

  Neighbor:  192.0.2.1
  Connection state:  bgp-st-idle

gRPC session timed out after 10 seconds: context deadline exceeded
```

## OpenConfig models

[http://ops.openconfig.net](http://ops.openconfig.net/branches/models/telemetry-version/)

- openconfig-acl
- openconfig-aft
- openconfig-bfd-ni
- openconfig-bfd
- openconfig-bgp-rib
- openconfig-bgp
- openconfig-catalog
- openconfig-interfaces
- openconfig-isis
- openconfig-lacp
- openconfig-lldp
- openconfig-local-routing
- openconfig-mpls
- openconfig-multicast
- openconfig-network-instance-sr-rsvp-coexistence
- openconfig-network-instance-sr
- openconfig-network-instance-srte-policy
- openconfig-network-instance
- openconfig-openflow
- openconfig-optical-amplifier
- openconfig-ospf
- openconfig-platform
- openconfig-probes
- openconfig-qos
- openconfig-relay-agent
- openconfig-routing-policy
- openconfig-stp
- openconfig-system
- openconfig-telemetry
- openconfig-terminal-device
- openconfig-transport-line-protection
- openconfig-types
- openconfig-vlan
- openconfig-wavelength-router

### Protobuf files

- [ems_grpc.proto](https://github.com/ios-xr/model-driven-telemetry/blob/master/protos/732/mdt_grpc_dialin/ems_grpc.proto)
- [telemetry.proto](https://github.com/ios-xr/model-driven-telemetry/blob/master/protos/732/telemetry.proto)
- [bgp_nbr_bag.proto](https://github.com/ios-xr/model-driven-telemetry/blob/master/protos/732/cisco_ios_xr_ipv4_bgp_oper/bgp/instances/instance/instance_active/default_vrf/afs/af/neighbor_af_table/neighbor/bgp_nbr_bag.proto)