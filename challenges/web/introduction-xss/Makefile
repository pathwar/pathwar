PREFIX ?= pathwar/
PATHWAR_OPTS ?=

.PHONY: pathwar.run
pathwar.run: pathwar.prepare
	pathwar $(PATHWAR_OPTS) compose up --force-recreate pathwar-compose.yml

.PHONY: pathwar.down
pathwar.down:
	pathwar --debug compose down $(notdir $(PWD))

.PHONY: pathwar.ps
pathwar.ps:
	pathwar compose ps | grep $(notdir $(PWD))

.PHONY: docker.build
docker.build:
	docker-compose build --pull

.PHONY: pathwar.prepare
pathwar.prepare:
	pathwar $(PATHWAR_OPTS) compose prepare --no-push . > pathwar-compose.yml

.PHONY: pathwar.push
pathwar.push:
	pathwar $(PATHWAR_OPTS) compose prepare --prefix=$(PREFIX) . > pathwar-compose.yml

.PHONY: docker.pathwar.push
docker.pathwar.push:
	docker run -it --rm --privileged -v /var/run/docker.sock:/var/run/docker.sock -v "$(PWD):$(PWD)" -w "$(PWD)" pathwar/agent compose prepare --prefix=$(PREFIX) . > pathwar-compose.yml

.PHONY: pathwar.register
pathwar.register:
	pathwar $(PATHWAR_OPTS) compose register --print ./pathwar-compose.yml > _register.sh
	chmod +x _register.sh

.PHONY: make.bump
make.bump:
	wget -O Makefile https://raw.githubusercontent.com/pathwar/pathwar/master/challenges/challenge.mk