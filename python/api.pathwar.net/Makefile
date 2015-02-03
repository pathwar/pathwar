BLUEPRINT_FILE ?=	apiary.apib
BLUEPRINT_TEMPLATE ?=	default
FIG_API_SERVICE ?=	api
NODE_SDK_SERVICE ?=	nodesdk


# FILES/DIRECTORIES
portal.pathwar.net:
	git clone https://github.com/pathwar/portal.pathwar.net


# ACTIONS
.PHONY:	all build release up shell clean kill stop

all:	build

build:	api_build blueprint_build

release:	gh-pages

up:	api_up portal_up
	fig logs

kill:
	fig kill

stop:
	fig stop

shell:	api_shell

clean:	blueprint_clean api_clean


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
.PHONY:	api_build api_up api_shell portal_up mongo_up smtp_up flush-db seed-db

api_build:	portal.pathwar.net
	fig build


flush-db:	mongo_up
	fig run --no-deps $(FIG_API_SERVICE) python run.py flush-db

seed-db:
	$(MAKE) seed-db-node

seed-db-node:	flush-db
	fig run --no-deps $(NODE_SDK_SERVICE) npm run seed

seed-db-python:	flush-db
	fig run --no-deps $(FIG_API_SERVICE) python run.py seed-db

api_up:
	fig kill api
	fig up --no-recreate -d api

api_shell:	mongo_up
	fig run --no-deps $(FIG_API_SERVICE) /bin/bash

mongo_up:
	fig up --no-recreate -d mongo

smtp_up:
	fig up --no-recreate -d smtp

portal_up:
	fig up --no-recreate -d portal

api_clean:
	fig stop
	fig rm --force
