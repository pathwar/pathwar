dev:
	docker build -t pathwar/level-training-http .
	pathwar.land hypervisor prune
	pathwar.land hypervisor run --web-port=8899 pathwar/level-training-http
