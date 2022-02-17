## Chapter 07 lab
07-up:
	cd ch07; sudo containerlab deploy -t ./topo.yml --reconfigure

## Chapter 07 cleanup
07-down:
	cd ch07; sudo containerlab destroy -t ./topo.yml --cleanup	

top-talkers-stop:
	sudo pkill -f ethr

# to run: cd ch07/top-talkers; sudo ip netns clab-netgo-host-2 go run main.go
top-talkers-start: 
	docker exec -d clab-netgo-cvx systemctl restart hsflowd
	docker exec -d clab-netgo-host-3 ./ethr -s
	docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.253 -b 900K -d 60s -p udp -l 1KB
	docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.252 -b 600K -d 60s -p udp -l 1KB
	docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.251 -b 400K -d 60s -p udp -l 1KB
	@echo "run 'cd ch07/top-talkers; sudo ip netns exec clab-netgo-host-2 go run main.go; cd ../../'"
	