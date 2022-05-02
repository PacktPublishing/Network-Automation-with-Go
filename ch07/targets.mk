## Build binary for CVX
build-go-cvx:
	cd ch07/ansible/cvx; go build -o ../library/go_cvx

## Build binary for SRL
build-go-srl:
	cd ch07/ansible/srl; go build -o ../library/go_srl

## Build binary for state validation
build-go-state:
	cd ch07/ansible/state; go build -o ../library/go_state