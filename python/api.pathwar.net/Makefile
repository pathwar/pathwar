CONFIG ?=		docker-compose-local.yml
DOCKER_COMPOSE ?=	docker-compose -p$(USER)_api -f$(CONFIG)


BLUEPRINT_FILE ?=	apiary.apib
BLUEPRINT_TEMPLATE ?=	default


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

REPOS ?=	portal.pathwar.net node-pathwar


all:	up


# FILES/DIRECTORIES
$(REPOS):
	git clone https://github.com/pathwar/$@


# ACTIONS
.PHONY:	all build release up shell clean kill stop pull re


re: stop up


pull:	$(REPOS)
	for repo in $(REPOS); do cd $$repo; git pull; cd -; done


build:	api_build blueprint_build


release:	gh-pages


up:	api_up portal_up


stop:
	$(DOCKER_COMPOSE) kill
	$(DOCKER_COMPOSE) stop
	$(DOCKER_COMPOSE) rm --force


logs kill ps:
	$(DOCKER_COMPOSE) $@


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
	$(DOCKER_COMPOSE) build


flush-db:	api_up
	$(DOCKER_COMPOSE) run --no-deps api python pathwar_api/run.py flush-db


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
	$(DOCKER_COMPOSE) run --no-deps nodesdk npm run seed


seed-db-python:	flush-db
	$(DOCKER_COMPOSE) run --no-deps api python pathwar_api/run.py seed-db


api_up:
	$(DOCKER_COMPOSE) kill api
	$(DOCKER_COMPOSE) rm --force api
	$(DOCKER_COMPOSE) up --no-recreate -d api


api_logs:
	$(DOCKER_COMPOSE) logs api


api_shell:	mongo_up
	$(DOCKER_COMPOSE) run --no-deps api /bin/bash


mongo_up:
	$(DOCKER_COMPOSE) up --no-recreate -d mongo


mongocli:
	$(DOCKER_COMPOSE) run mongo /bin/sh -ec 'mongo $${MONGO_HOST:-$${MONGO_PORT_27017_TCP_ADDR}}/pathwar'


smtp_up:
	$(DOCKER_COMPOSE) up --no-recreate -d smtp


portal_up:
	$(DOCKER_COMPOSE) kill portal
	$(DOCKER_COMPOSE) rm --force portal
	$(DOCKER_COMPOSE) up --no-recreate -d portal


portal_clean:
	-$(DOCKER_COMPOSE) kill portal
	-$(DOCKER_COMPOSE) rm --force portal
	-rm -rf portal.pathwar.net/build


api_clean:
	-$(DOCKER_COMPOSE) stop api
	-$(DOCKER_COMPOSE) rm --force api
