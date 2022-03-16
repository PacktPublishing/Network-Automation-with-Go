## Base topology lab
lab-up:
	sudo containerlab deploy -t topo-base/topo.yml --reconfigure

## Base topology cleanup
lab-down:
	sudo containerlab destroy -t topo-base/topo.yml --cleanup	