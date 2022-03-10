## Base lab
topo-up: ## Bring up base lab
	cd topo; sudo containerlab deploy -t ./topo.yml --reconfigure
	cd -

## Base lab cleanup
topo-down: ## Turn down base lab
	cd topo; sudo containerlab destroy -t ./topo.yml --cleanup
	cd -

.PHONY: host
host: ## Build host image
	cd topo; docker build -t thebook:host host/ -f host/Dockerfile