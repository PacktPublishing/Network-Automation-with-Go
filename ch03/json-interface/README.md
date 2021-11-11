# JSON 

## Decoding

### Input
```json
{
  "router": [
    {
      "hostname": "router1.example.com",
      "ip": "192.0.2.1",
      "asn": 64512
    },
    {
      "hostname": "router2.example.com",
      "ip": "198.51.100.1",
      "asn": 65535
    }
  ]
}
```

### Code
Decode reads the next JSON-encoded value from its input and stores it in the empty interface

```go
func main() {
	file, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	d := json.NewDecoder(file)

	var empty map[string]interface{}

	err = d.Decode(&empty)
	if err != nil {
		panic(err)
	}

	for _, r := range empty["router"].([]interface{}) {
		fmt.Printf("%v\n", r)
	}

}

```

### Output
```bash
⇨  go run main.go 
map[asn:64512 hostname:router1.example.com ip:192.0.2.1]
map[asn:65535 hostname:router2.example.com ip:198.51.100.1]
```
