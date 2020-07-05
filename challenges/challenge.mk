PREFIX ?= pathwar/
PATHWAR_OPTS ?=

.PHONY: docker.build
docker.build:
	docker-compose build --pull

.PHONY: pathwar-prepare
pathwar.prepare:
	pathwar $(PATHWAR_OPTS) compose prepare --prefix=$(PREFIX) --no-push .

.PHONY: pathwar-push
pathwar.push:
	pathwar $(PATHWAR_OPTS) compose prepare --prefix=$(PREFIX) .
