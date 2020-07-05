PREFIX ?= pathwar/
PATHWAR_OPTS ?=

.PHONY: pathwar.run
pathwar.run:
	pathwar $(PATHWAR_OPTS) compose prepare --no-push ./ > pathwar-compose.yml
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

.PHONY: pathwar-prepare
pathwar.prepare:
	pathwar $(PATHWAR_OPTS) compose prepare --prefix=$(PREFIX) --no-push .

.PHONY: pathwar-push
pathwar.push:
	pathwar $(PATHWAR_OPTS) compose prepare --prefix=$(PREFIX) .

.PHONY: make.bump
make.bump:
	wget -O rules.mk https://raw.githubusercontent.com/tree/master/challenges/challenge.mk Makefile
