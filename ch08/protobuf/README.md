## Protocol Buffers

Requirements:

Redhad:
```bash
sudo yum install protobuf-compiler
```

Ubuntu:
```bash
 sudo apt install protobuf-compiler
```

```
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
$ cd write && go run protobuf
```

### Read

```bash
$ cd read && go run protobuf
uplinks:{name:"Ethernet1"  prefix:"192.0.2.1/31"}  uplinks:{name:"Ethernet2"  prefix:"192.0.2.2/31"}  peers:{ip:"192.0.2.0"  asn:65000}  peers:{ip:"192.0.2.3"  asn:65002}  asn:65001  loopback:{ip:"198.51.100.1"}
```

### Inspecting the data

1. You can see the content with the `hexdump` command.

```hexdump
$ hexdump -C router.data
00000000  0a 19 0a 09 45 74 68 65  72 6e 65 74 31 12 0c 31  |....Ethernet1..1|
00000010  39 32 2e 30 2e 32 2e 31  2f 33 31 0a 19 0a 09 45  |92.0.2.1/31....E|
00000020  74 68 65 72 6e 65 74 32  12 0c 31 39 32 2e 30 2e  |thernet2..192.0.|
00000030  32 2e 32 2f 33 31 12 0f  0a 09 31 39 32 2e 30 2e  |2.2/31....192.0.|
00000040  32 2e 30 10 e8 fb 03 12  0f 0a 09 31 39 32 2e 30  |2.0........192.0|
00000050  2e 32 2e 33 10 ea fb 03  18 e9 fb 03 22 0e 0a 0c  |.2.3........"...|
00000060  31 39 38 2e 35 31 2e 31  30 30 2e 31              |198.51.100.1|
```

2. You could also out the protobuf encoded slice of bytes as follows

```go
	out, err := proto.Marshal(router)
	// parse error
	fmt.Printf("%X", out)
```

3. If we group the byte content for convenience, we get something like:

```hexdump
0A 19 0A 09 45 74 68 65 72 6E 65 74 31 
12 0C 31 39 32 2E 30 2E 32 2E 31 2F 33 31 
0A 19 0A 09 45 74 68 65 72 6E 65 74 32 
12 0C 31 39 32 2E 30 2E 32 2E 32 2F 33 31 
12 0F 0A 09 31 39 32 2E 30 2E 32 2E 30 10 E8 FB 03 
12 0F 0A 09 31 39 32 2E 30 2E 32 2E 33 10 EA FB 03 
18 E9 FB 03 22 
0E 0A 0C 31 39 38 2E 35 31 2E 31 30 30 2E 31
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
E9:
FB:
03:
22  tag: loopback(4), field encoding: LENGTH_DELIMITED(2)
0E  "loopback".length(): 14
0A  tag: ip(1), field encoding: LENGTH_DELIMITED(2)
0C  "ip".length(): 12 
...
```

#### Varint

- Most significant bit (msb) set – this indicates that there are further bytes to come
- The lower 7 bits of each byte are used to store the two's complement representation of the number in groups of 7 bits

```
18 tag: asn(3), field encoding: VARINT(0)
E9: 1 1101001
FB: 1 1111011
03: 0 0000011
```

- Reverse the two groups of 7 bits because

```
0000011 1111011 1101001
```

- Concatenate them to get your final value

```
000 0011 ++ 111 1011 ++ 110 1001
11 1111011 1101001
(1111110111101001)₂
```

- In decimal

```
(1 × 2¹⁵)+(1 × 2¹⁴)+(1 × 2¹³)+(1 × 2¹²)+(1 × 2¹¹)+(1 × 2¹⁰)+(0 × 2⁹)+(1 × 2⁸)+(1 × 2⁷)+(1 × 2⁶)+(1 × 2⁵)+(0 × 2⁴)+(1 × 2³)+(0 × 2²)+(0 × 2¹)+(1 × 2⁰)
= 32768 + 16384 + 8192 + 4096 + 2048 + 1024 + 256 + 128 + 64 + 32 + 8 + 1
= (65001)₁₀
```

4. Alternatively, `protoc` can show you the decoded data.

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

### Size

```bash
$ ls -ls router* | awk '{print $6, $10}'
108 router.data
454 router_indent.json
220 router.json
```