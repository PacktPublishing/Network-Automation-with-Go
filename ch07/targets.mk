## Build binary for CVX
build-go-cvx:
	cd ch07/cvx; go build -o ../library/go_cvx

## Build binary for SRL
build-go-srl:
	cd ch07/srl; go build -o ../library/go_srl