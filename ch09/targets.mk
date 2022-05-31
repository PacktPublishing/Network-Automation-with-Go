
gnmic-start: 
	cd ch09/; docker-compose up -d; cd ../

bgp-ping-stop:
	cd ch10/bgp-ping; docker-compose down; cd ../../
