## Base topology lab
lab-base:
	sudo containerlab deploy -t topo-base/topo.yml --reconfigure

## Base topology cleanup
base-down:
	sudo containerlab destroy -t topo-base/topo.yml --cleanup	