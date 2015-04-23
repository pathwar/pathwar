SERVICES ?=	$(shell cat docker-compose.yml | grep '^[a-z]' | cut -d: -f1)
# By default $(SERVICE) is the first service in docker-compose.yml
SERVICE ?=	$(shell echo $(SERVICES) | tr " " "\n" | head -1)

PACKAGE_NAME = package-$(shell basename $(shell pwd)).tar

## Actions
.PHONY: all build run shell $(PACKAGE_NAME)

all:	build

build:
	docker-compose build

run:
	docker-compose stop
	docker-compose up -d
	docker-compose ps
	docker-compose logs

shell:
	docker-compose run $(SERVICE) /bin/bash


## Package
package-level:
	curl -O https://raw.githubusercontent.com/pathwar/core/master/mk/package-level
	chmod +x $@

package: $(PACKAGE_NAME)

$(PACKAGE_NAME): package-level
	-rm -f $(PACKAGE_NAME)
	./package-level build


## Travis
.PHONY: travis_install travis_run travis_run_service

travis_install:
	# Install travis-docker
	curl -sLo - https://github.com/moul/travis-docker/raw/master/install.sh | sh -xe

travis_run:
	for service in $(SERVICES); do \
	  $(MAKE) travis_run_service SERVICE=$$service; \
	done

travis_run_service:
	# ./run is a wrapper to allow travis to run docker commands
	./run docker-compose run $(SERVICE) /bin/bash -d 'echo OK'
