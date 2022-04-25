## Nautobot

PoC

```bash
cd ch06/nautobot/client
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
wget https://develop.demo.nautobot.com/api/swagger.yaml\?api_version\=1.3 -O swagger.yaml
oapi-codegen -generate client -o nautobot.go -package nautobot swagger.yaml
oapi-codegen -generate types -o types.go -package nautobot swagger.yaml
go mod init github.com/nautobot/go-nautobot
```

```bash
cd ..
go mod edit -replace github.com/nautobot/go-nautobot=./client
```