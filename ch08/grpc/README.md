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

## Running the example

By default the example runs in the `gpb` subscription mode. To run in the self-describing kv mode add the `-kvmode=true` flag.

```bash
$ go run grpc

BGP config applied on sandbox-iosxr-1.cisco.com:57777


Streaming telemetry from sandbox-iosxr-1.cisco.com:57777

----
Time: Tue May 17 13:24:35 2022
Path: Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor

  Neighbor:  192.0.2.1
  Connection state:  bgp-st-idle

----
Time: Tue May 17 13:24:37 2022
Path: Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor

  Neighbor:  192.0.2.1
  Connection state:  bgp-st-idle

----
Time: Tue May 17 13:24:39 2022
Path: Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor

  Neighbor:  192.0.2.1
  Connection state:  bgp-st-idle

gRPC session timed out after 10 seconds: context deadline exceeded
```

## OpenConfig models

- [http://ops.openconfig.net](http://ops.openconfig.net/branches/models/telemetry-version/)


## Protobuf files

- [ems_grpc.proto](https://github.com/ios-xr/model-driven-telemetry/blob/master/protos/732/mdt_grpc_dialin/ems_grpc.proto)
- [telemetry.proto](https://github.com/ios-xr/model-driven-telemetry/blob/master/protos/732/telemetry.proto)
