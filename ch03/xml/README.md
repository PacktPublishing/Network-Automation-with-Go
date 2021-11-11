# XML

## Decoding

### Input
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<routers>
  <router>
    <hostname>router1.example.com</hostname>
    <ip>192.0.2.1</ip>
    <asn>64512</asn>
  </router>
  <router>
    <hostname>router2.example.com</hostname>
    <ip>198.51.100.1</ip>
    <asn>65535</asn>
  </router>
</routers>
```

### Code
Decode reads the next XML-encoded value from its input and stores it in the value pointed to by v.

```go
func main() {
	file, err := os.Open("input.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	src := xml.NewDecoder(file)

	var inv Inventory
	err = src.Decode(&inv)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", inv)
}
```

### Output
```bash
⇨  go run main.go 
{Routers:[{Hostname:router1.example.com IP:192.0.2.1 ASN:64512} {Hostname:router2.example.com IP:198.51.100.1 ASN:65535}]}
```