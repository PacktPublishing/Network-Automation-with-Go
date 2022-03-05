## Chapter 07 lab
07-up:
	cd ch07; sudo containerlab deploy -t ./topo.yml --reconfigure

## Chapter 07 cleanup
07-down:
	cd ch07; sudo containerlab destroy -t ./topo.yml --cleanup	

bgp-ping-build:
	cd ch07/bgp-ping; go build -o bgp-ping main.go

bgp-ping-stop:
	cd ch07/bgp-ping; docker-compose down; cd ../../

bgp-ping-start: bgp-ping-build
	docker exec -d clab-netgo-host-3 /workdir/bgp-ping/bgp-ping -id host-3 -nlri 100.64.0.2 -laddr 203.0.113.254 -raddr 203.0.113.129 -las 65005 -ras 65002 -p
	docker exec -d clab-netgo-host-1 /workdir/bgp-ping/bgp-ping -id host-1 -nlri 100.64.0.0 -laddr 203.0.113.0 -raddr 203.0.113.1 -las 65003 -ras 65000 -p
	docker exec -d clab-netgo-host-2 /cloudprober -config_file /workdir/workdir/cloudprober.cfg
	cd ch07/bgp-ping; docker-compose up -d; cd ../../
	@echo 'http://localhost:3000'

traffic-start:
	docker exec -d clab-netgo-cvx systemctl restart hsflowd
	docker exec -d clab-netgo-host-3 ./ethr -s
	docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.253 -b 900K -d 60s -p udp -l 1KB
	docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.252 -b 600K -d 60s -p udp -l 1KB
	docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.251 -b 400K -d 60s -p udp -l 1KB

top-talkers-start: traffic-start
	@echo "run 'cd ch07/top-talkers; sudo ip netns exec clab-netgo-host-2 go run main.go; cd ../../'"
	
07-stop: bgp-ping-stop
	sudo pkill -f bgp-ping
	sudo pkill -f cloudprober
	sudo pkill -f ethr


capture-start: traffic-start
	cd ch07/packet-capture; go build -o packet-capture main.go
	docker exec -it clab-netgo-host-2 /workdir/packet-capture/packet-capture

capture-debug:
	echo "docker exec -it clab-netgo-host-2 bash -c 'cd /workdir/packet-capture/; dlv debug main.go --listen=:2345 --headless --api-version=2'"

