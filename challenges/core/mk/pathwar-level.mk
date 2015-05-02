SERVICES ?=	$(shell cat docker-compose.yml | grep '^[a-z]' | cut -d: -f1)
# By default $(SERVICE) is the first service in docker-compose.yml
SERVICE ?=	$(shell echo $(SERVICES) | tr " " "\n" | head -1)
S3_URL ?= 	s3://pathwar-levels/
PACKAGE_NAME ?=	package-$(shell basename $(shell pwd)).tar
EXCLUDES ?=	$(PACKAGE_NAME) .git
ASSETS ?=	$(filter-out $(EXCLUDES),$(wildcard *))
PKGLVL ?=	/tmp/package_level


## Actions
all: build
.PHONY: all build run shell package publish_on_s3 info


info:
	@echo "ASSETS:        $(ASSETS)"
	@echo "EXLUDES:       $(EXCLUDES)"
	@echo "PACKAGE_NAME:  $(PACKAGE_NAME)"
	@echo "S3_URL:        $(S3_URL)"
	@echo "SERVICE:       $(SERVICE)"
	@echo "SERVICES:      $(SERVICES)"


build:
	docker-compose build


run:
	docker-compose kill
	docker-compose stop
	docker-compose up -d
	docker-compose ps
	docker-compose logs


shell:
	docker-compose run $(SERVICE) /bin/bash


## Package
$(PKGLVL):
	curl -s https://raw.githubusercontent.com/pathwar/core/master/mk/package-level -o $@
	chmod +x $@


package: $(PACKAGE_NAME)


$(PACKAGE_NAME): $(PKGLVL) $(ASSETS)
	$(PKGLVL) build


publish_on_s3: $(PACKAGE)
	s3cmd put --acl-public $(PACKAGE_NAME) $(S3_URL)/$(shell head -c 128 /dev/urandom | tr -dc A-Za-z0-9).tar


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
