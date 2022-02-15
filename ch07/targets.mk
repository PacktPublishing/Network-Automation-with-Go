## Chapter 07 lab
07-up:
	cd ch07; sudo containerlab deploy -t ./topo.yml --reconfigure

## Chapter 07 cleanup
07-down:
	cd ch07; sudo containerlab destroy -t ./topo.yml --cleanup	