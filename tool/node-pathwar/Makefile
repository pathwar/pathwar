.PHONY:	build test watch_test

build:
	npm build


test:	node_modules
	npm test


watch_test:	node_modules
	while true; do clear; npm test; sleep .5; fswatch -1 *; done


node_modules:
	npm install


docker_test:
	fig up -d --no-recreate api mongo
	fig run sdk npm test


watch_docker_test:
	fig up -d --no-recreate api mongo
	while true; do clear; fig run sdk npm test; sleep .5; fswatch -1 *; done
