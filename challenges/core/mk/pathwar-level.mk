CONFIG ?=		docker-compose.yml
SERVICES ?=		$(shell cat docker-compose.yml | grep '^[a-z]' | cut -d: -f1)
# By default $(SERVICE) is the first service in docker-compose.yml
SERVICE ?=		$(shell echo $(SERVICES) | tr " " "\n" | head -1)
S3_URL ?=		s3://pathwar-levels/
PACKAGE_NAME ?=		package-$(shell basename $(shell pwd)).tar
EXCLUDES ?=		$(PACKAGE_NAME) .git
ASSETS ?=		$(filter-out $(EXCLUDES),$(wildcard *))
PKGLVL ?=		/tmp/package_level
DOCKER_COMPOSE_NAME ?=	$(USER)$(shell basename $(shell pwd) | sed 's/[^a-z]//g')
DOCKER_COMPOSE ?=	docker-compose -p$(DOCKER_COMPOSE_NAME)
STORE_HOST ?=		store.pathwar.net
STORE_PATH ?=		pathwar
S3CMD ?=		s3cmd
SECTIONS ?=		$(shell cat $(CONFIG) | grep -E '^[a-z]' | cut -d: -f1)
MAIN_SECTION ?=		$(shell echo $(SECTIONS) | cut -d\  -f1)
MAIN_CID :=		$(shell $(DOCKER_COMPOSE) ps -q $(MAIN_SECTION))
EXEC_MAIN_SECTION :=	docker exec -it $(MAIN_CID)
UNIX_USER ?=		bobby


## Actions
all: up ps logs
.PHONY: all build run shell package info


shellmysql:
	docker run \
	  -it --rm \
	  orchardup/mysql \
	  mysql \
	    -h$(shell docker inspect -f '{{ .NetworkSettings.IPAddress }}' $(shell $(DOCKER_COMPOSE) ps -q mysql))


unix_run: build
	docker run $(DOCKER_COMPOSE_NAME)_$(MAIN_SECTION)
	docker commit `docker ps -lq` tmp
	docker run -it --rm -u $(UNIX_USER) tmp level-enter


info: before
	@echo "ASSETS:        $(ASSETS)"
	@echo "EXLUDES:       $(EXCLUDES)"
	@echo "PACKAGE_NAME:  $(PACKAGE_NAME)"
	@echo "S3_URL:        $(S3_URL)"
	@echo "SERVICE:       $(SERVICE)"
	@echo "SERVICES:      $(SERVICES)"


before::


kill stop rm:: before
	$(DOCKER_COMPOSE) kill $(SECTIONS)
	$(DOCKER_COMPOSE) stop $(SECTIONS)
	$(DOCKER_COMPOSE) rm --force $(SECTIONS)


re:: stop up


up:: before
	$(DOCKER_COMPOSE) up --no-recreate -d $(SECTIONS)


shell:
	$(DOCKER_COMPOSE) run $(SERVICE) /bin/bash


shellexec:: before
	$(EXEC_MAIN_SECTION) bash


## Package
$(PKGLVL):
	curl -s https://raw.githubusercontent.com/pathwar/core/master/mk/package-level -o $@
	chmod +x $@


package: $(PACKAGE_NAME)


$(PACKAGE_NAME): $(PKGLVL) $(ASSETS)
	-@docker ps -a -f name=pathwar-exportme -q | xargs docker rm 2>/dev/null || true
	$(PKGLVL) build


build ps:: before
	$(DOCKER_COMPOSE) $@


logs:: up
	$(DOCKER_COMPOSE) $@ $(SECTIONS)


.PHONY: publish_on_store
publish_on_store: $(PACKAGE_NAME)
	$(eval RAND := $(shell openssl rand -base64 46 | tr -dc A-Za-z0-9))
	rsync -Pave ssh $(PACKAGE_NAME) $(STORE_HOST):store/$(STORE_PATH)/$(RAND).tar
	@echo http://$(STORE_HOST)/$(STORE_PATH)/$(RAND).tar


.PHONY: publish_on_s3
publish_on_s3: $(PACKAGE_NAME)
	$(eval RAND := $(shell openssl rand -base64 46 | tr -dc A-Za-z0-9))
	$(S3CMD) put --acl-public $(PACKAGE_NAME) $(S3_URL)/$(RAND).tar
	$(S3CMD) info $(S3_URL)/$(RAND).tar | grep URL | awk '{print $$2}'


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
	./run "$(DOCKER_COMPOSE) run $(SERVICE) /bin/bash -xec 'echo OK'"
