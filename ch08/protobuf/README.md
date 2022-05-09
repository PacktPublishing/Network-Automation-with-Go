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

->

```go
type Router struct {
	Uplink   []*Uplink `protobuf:"bytes,1,rep,name=uplink,proto3" json:"uplink,omitempty"`
	Peer     []*Peer   `protobuf:"bytes,2,rep,name=peer,proto3" json:"peer,omitempty"`
	Asn      int32     `protobuf:"varint,3,opt,name=asn,proto3" json:"asn,omitempty"`
	Loopback *Addr     `protobuf:"bytes,4,opt,name=loopback,proto3" json:"loopback,omitempty"`
}
```