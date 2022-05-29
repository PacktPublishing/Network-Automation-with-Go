module nautobot

go 1.17

require (
	github.com/deepmap/oapi-codegen v1.11.0
	github.com/nautobot/go-nautobot v0.0.0-00010101000000-000000000000
)

require github.com/google/uuid v1.3.0 // indirect

replace github.com/nautobot/go-nautobot => ./client
