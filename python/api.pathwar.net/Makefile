BLUEPRINT_FILE ?=	apiary.apib
BLUEPRINT_TEMPLATE ?=	default
DOCKER_COMPOSE_FILE ?=	docker-compose.yml

WATCH_FILES ?=	$(shell \
		  find node-pathwar/ ./ \
		  -not -ipath '*/node_modules/*' \
		  -not -ipath '*/.git/*' \
		  -not -ipath '*/vendor/*' \
		  -not -ipath '*~' \
		  -not -ipath '*.pyc' \
		  -not -ipath '*\#' \
		  -type f \
		)


# FILES/DIRECTORIES
portal.pathwar.net:
	git clone https://github.com/pathwar/portal.pathwar.net


# ACTIONS
.PHONY:	all build release up shell clean kill stop

all:	build

build:	api_build blueprint_build

release:	gh-pages

up:	api_up portal_up
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs

kill:
	docker-compose -f $(DOCKER_COMPOSE_FILE) kill

stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop

shell:	api_shell

clean:	blueprint_clean api_clean portal_clean


# BLUEPRINT
.PHONY:	blueprint_watch blueprint_build blueprint_clean

blueprint_watch:
	aglio -i $(BLUEPRINT_FILE) -t $(BLUEPRINT_TEMPLATE) -s

blueprint_build:
	aglio -i $(BLUEPRINT_FILE) -t $(BLUEPRINT_TEMPLATE) -o apiary.html

blueprint_clean:
	-rm -f apiary.html index.html


# GH-PAGES
.PHONY:	gh-pages gh-pages_do gh-pages_teardown

gh-pages:
	$(MAKE) gh-pages_do || $(MAKE) gh-pages_teardown

gh-pages_do:
	git branch -D gh-pages || true
	git checkout -b gh-pages
	$(MAKE) blueprint_build
	mv apiary.html index.html
	git add index.html
	git commit index.html -m "Rebuild assets"
	git push -u origin gh-pages -f
	$(MAKE) gh-pages_teardown

gh-pages_teardown:
	git checkout master


# TRAVIS
.PHONY:	travis

travis:
	find . -name Dockerfile | xargs cat | grep -vi ^maintainer | bash -n
	aglio -i $(BLUEPRINT_FILE) -t $(BLUEPRINT_TEMPLATE) -o apiary.html


# API
.PHONY:	api_build api_up api_shell portal_up mongo_up smtp_up flush-db seed-db portal_clean

api_build:	portal.pathwar.net node-pathwar
	docker-compose -f $(DOCKER_COMPOSE_FILE) build


flush-db:	mongo_up
	docker-compose -f $(DOCKER_COMPOSE_FILE) run --no-deps api python pathwar_api/run.py flush-db

seed-db:
	$(MAKE) seed-db-node

seed-db-watch:
	while true; do \
	  clear; \
	  $(MAKE) seed-db; \
	  sleep .5; \
	  fswatch -1 $(WATCH_FILES); \
	done

seed-db-node:	flush-db node-pathwar
	docker-compose -f $(DOCKER_COMPOSE_FILE) run --no-deps nodesdk npm run seed

seed-db-python:	flush-db
	docker-compose -f $(DOCKER_COMPOSE_FILE) run --no-deps api python pathwar_api/run.py seed-db

api_up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) kill api
	docker-compose -f $(DOCKER_COMPOSE_FILE) rm --force api
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --no-recreate -d api

api_shell:	mongo_up
	docker-compose -f $(DOCKER_COMPOSE_FILE) run --no-deps api /bin/bash

mongo_up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --no-recreate -d mongo

smtp_up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --no-recreate -d smtp

portal_up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) kill portal
	docker-compose -f $(DOCKER_COMPOSE_FILE) rm --force portal
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --no-recreate -d portal


portal_clean:
	-docker-compose -f $(DOCKER_COMPOSE_FILE) kill portal
	-docker-compose -f $(DOCKER_COMPOSE_FILE) rm --force portal
	-rm -rf portal.pathwar.net/build

api_clean:
	-docker-compose -f $(DOCKER_COMPOSE_FILE) stop api
	-docker-compose -f $(DOCKER_COMPOSE_FILE) rm --force api

node-pathwar:
	git clone git://github.com/pathwar/node-pathwar
