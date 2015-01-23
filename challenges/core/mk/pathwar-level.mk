SERVICES ?=	$(shell cat fig.yml | grep '^[a-z]' | cut -d: -f1)
# By default $(SERVICE) is the first service in fig.yml
SERVICE ?=	$(shell echo $(SERVICES) | tr " " "\n" | head -1)


## Actions
.PHONY: all build run shell

all:	build

build:
	fig build

run:
	fig stop
	fig up -d
	fig ps
	fig logs

shell:
	fig run $(SERVICE) /bin/bash


## Travis
.PHONY: travis_install travis_run travis_run_service

travis_install:
	# Install travis-docker
	curl -sLo - https://github.com/moul/travis-docker/raw/master/install.sh | sh -xe

travis_run:
	for service in $(SERVICES); do \
	  $(MAKE) travis_run_service SERVICE=$$servie; \
	done

travis_run_service:
	# ./run is a wrapper to allow travis to run docker commands
	./run fig run $(SERVICE) /bin/bash -d 'echo OK'
