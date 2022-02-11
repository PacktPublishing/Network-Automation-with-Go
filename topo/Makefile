
up: 
	sudo containerlab deploy -t topo.yml --reconfigure

down:
	sudo containerlab destroy -t topo.yml

.PHONY: host
host:
	docker build -t thebook:host host/ -f host/Dockerfile