
gnmic-start: 
	docker exec -d clab-netgo-host-3 \
	gnmic --config /workdir/topo-full/workdir/gnmic.yaml \
	subscribe
	cd ch09/; docker-compose up -d; cd ../
	@echo "run 'sudo ip netns exec clab-netgo-host-2 ${GOBIN} run ch09/event-manager/main.go'"

gnmic-stop: 
	cd ch09/; docker-compose down; cd ../
	sudo pkill -f gnmic

