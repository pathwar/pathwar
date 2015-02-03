.PHONY:	build test watch_test fig_api


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


docker_test:	fig_api
	fig run sdk npm run seed
	fig run sdk npm run test


watch_docker_test:	fig_api
	while true; \
	  do clear; \
	  fig run sdk npm test; \
	  sleep .5; \
	  fswatch -1 *; \
	done


fig_api:	api.pathwar.net
	fig up -d --no-recreate api
	fig run --no-deps api python run.py flush-db


api.pathwar.net:
	git clone https://github.com/pathwar/api.pathwar.net
