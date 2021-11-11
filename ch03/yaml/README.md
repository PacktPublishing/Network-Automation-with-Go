# YAML

## Decoding

### Input
```yaml
router:
  - hostname: "router1.example.com"
    ip: "192.0.2.1"
    asn: 64512
  - hostname: "router2.example.com"
    ip: "198.51.100.1"
    asn: 65535
```

### Code
Decode reads the next YAML-encoded value from its input and stores it in the value pointed to by v.

```go
func main() {
	file, err := os.Open("input.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	src := yaml.NewDecoder(file)
	
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
go run main.go 
{Routers:[{Hostname:router1.example.com IP:192.0.2.1 ASN:64512} {Hostname:router2.example.com IP:198.51.100.1 ASN:65535}]}
```