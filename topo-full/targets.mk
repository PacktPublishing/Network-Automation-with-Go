## Full topology lab
lab-full:
	sudo containerlab deploy -t topo-full/topo.yml --reconfigure

## Full topology cleanup
full-down:
	sudo containerlab destroy -t topo-full/topo.yml --cleanup	