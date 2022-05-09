## Protocol Buffers

Requirements:

```bash
sudo yum install protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

### Creating bindings 

```bash
protoc --go_out=. model.proto
```

=>

```go
type Router struct {
	Uplinks  []*Uplink `protobuf:"bytes,1,rep,name=uplinks,proto3" json:"uplinks,omitempty"`
	Peers    []*Peer   `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty"`
	Asn      int32     `protobuf:"varint,3,opt,name=asn,proto3" json:"asn,omitempty"`
	Loopback *Addr     `protobuf:"bytes,4,opt,name=loopback,proto3" json:"loopback,omitempty"`
}
```