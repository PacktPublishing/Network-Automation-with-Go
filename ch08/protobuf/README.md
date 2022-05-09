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

```hexdump
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

Let's print out the protobuf encoded slice of bytes

```go
	out, err := proto.Marshal(router)
	// parse error
	fmt.Printf("%X", out)
```

After grouping the output for convenience, we get something like:

```hexdump
0A 19 0A 09 45 74 68 65 72 6E 65 74 31 
12 0C 31 39 32 2E 30 2E 32 2E 31 2F 33 31 
0A 19 0A 09 45 74 68 65 72 6E 65 74 32 
12 0C 31 39 32 2E 30 2E 32 2E 32 2F 33 31 
12 0F 0A 09 31 39 32 2E 30 2E 32 2E 30 10 E8 FB 03 
12 0F 0A 09 31 39 32 2E 30 2E 32 2E 33 10 EA FB 03 
18 E9 FB 03
22 0E 0A 0C 31 39 38 2E 35 31 2E 31 30 30 2E 31
```

Protobuf uses Varint to serialize integers. The last three bits of the number store the wire type (field encoding), then right-shift by three to get the field number. Having this in mind and how to convert Hex to ASCII, this translates to:

```bash
Hex  Description
0A  tag: uplinks(1), field encoding: LENGTH_DELIMITED(2)
19  "uplinks".length(): 25
0A  tag: name(1), field encoding: LENGTH_DELIMITED(2)
09  "name".length(): 9 
45 'e'
74 't'
75 'h'
65 'e'
72 'r'
6E 'n'
65 'e'
74 't'
31 '1'
12 tag: prefix(2), field encoding: LENGTH_DELIMITED(2)
0C "prefix".length(): 12
31 '1'
39 '9'
32 '2'
2E '.'
...
2F '/'
33 '3'
31 '1'
...
12 tag: peers(2), field encoding: LENGTH_DELIMITED(2)
0F "prefix".length(): 15
...
18 tag: asn(3), field encoding: VARINT(0)
...
22  tag: loopback(4), field encoding: LENGTH_DELIMITED(2)
0E  "loopback".length(): 14
0A  tag: ip(1), field encoding: LENGTH_DELIMITED(2)
0C  "name".length(): 12 
...
```

3.

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

=>

```json
1 {
  1: "Ethernet1"
  2: "192.0.2.1/31"
}
```

### Size

```bash
$ ls -ls router* | awk '{print $6, $10}'
108 router.data
454 router_indent.json
220 router.json
```