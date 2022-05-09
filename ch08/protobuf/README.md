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

### Write

```bash
cd write && go run main.go
```

### Read

```bash
cd read && go run main.go
```

### Inspecting the data

1.

```bash
$ hexdump -c router.data
0000000  \n 031  \n  \t   E   t   h   e   r   n   e   t   1 022  \f   1
0000010   9   2   .   0   .   2   .   1   /   3   1  \n 031  \n  \t   E
0000020   t   h   e   r   n   e   t   2 022  \f   1   9   2   .   0   .
0000030   2   .   2   /   3   1 022 017  \n  \t   1   9   2   .   0   .
0000040   2   .   0 020 350 373 003 022 017  \n  \t   1   9   2   .   0
0000050   .   2   .   3 020 352 373 003 030 351 373 003   " 016  \n  \f
0000060   1   9   8   .   5   1   .   1   0   0   .   1                
000006c
```

2.

```bash
cat router.data | protoc --decode_raw
1 {
  1: "Ethernet1"
  2: "192.0.2.1/31"
}
1 {
  1: "Ethernet2"
  2: "192.0.2.2/31"
}
2 {
  1 {
    6: 0x302e322e302e3239
  }
  2: 65000
}
2 {
  1 {
    6: 0x332e322e302e3239
  }
  2: 65002
}
3: 65001
4 {
  1: "198.51.100.1"
}
```

