.PHONY:	build test watch_test docker_api


build:
	npm build


test:	node_modules
	npm test


watch_test:	node_modules
	while true; \
	  do clear; \
	  npm test; \
	  sleep .5; \
	  fswatch -1 *; \
	done


node_modules:
	npm install


docker_test:	docker_api
	docker-compose run sdk npm run seed
	docker-compose run sdk npm run test


watch_docker_test:	docker_api
	while true; \
	  do clear; \
	  docker-compose run sdk npm test; \
	  sleep .5; \
	  fswatch -1 *; \
	done


docker_api:	api.pathwar.net
	docker-compose up -d --no-recreate api
	docker-compose run --no-deps api python run.py flush-db


api.pathwar.net:
	git clone https://github.com/pathwar/api.pathwar.net
