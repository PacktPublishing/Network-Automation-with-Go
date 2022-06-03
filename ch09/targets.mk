
gnmic-start: 
	docker exec -d clab-netgo-host-3 \
	gnmic --config /workdir/topo-full/workdir/gnmic.yaml \
	subscribe
	cd ch09/; docker-compose up -d; cd ../
	@echo "run 'cd ch09/event-trigger; sudo ip netns exec clab-netgo-cvx ${GOBIN} run main.go; cd ../../'"

gnmic-stop: 
	cd ch09/; docker-compose down; cd ../
	sudo pkill -f gnmic

