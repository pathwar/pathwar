CONFIG ?=	docker-compose.yml
SERVICES ?=	$(shell cat docker-compose.yml | grep '^[a-z]' | cut -d: -f1)
# By default $(SERVICE) is the first service in docker-compose.yml
SERVICE ?=	$(shell echo $(SERVICES) | tr " " "\n" | head -1)
S3_URL ?= 	s3://pathwar-levels/
PACKAGE_NAME ?=	package-$(shell basename $(shell pwd)).tar
EXCLUDES ?=	$(PACKAGE_NAME) .git
ASSETS ?=	$(filter-out $(EXCLUDES),$(wildcard *))
PKGLVL ?=	/tmp/package_level
DOCKER_COMPOSE ?=	docker-compose -p$(USER)_$(shell basename $(shell pwd) | sed 's/[^a-z]//g')
S3CMD ?=	s3cmd
SECTIONS ?=	$(shell cat $(CONFIG) | grep -E '^[a-z]' | cut -d: -f1)
MAIN_SECTION ?=	$(shell echo $(SECTIONS) | cut -d\  -f1)
MAIN_CID :=	$(shell $(DOCKER_COMPOSE) ps -q $(MAIN_SECTION))
EXEC_MAIN_SECTION :=	docker exec -it $(MAIN_CID)
MYSQL_CID :=	$(shell $(DOCKER_COMPOSE) ps -q mysql)
MYSQL_IP :=	$(shell docker inspect -f '{{ .NetworkSettings.IPAddress }}' $(MYSQL_CID))


## Actions
all: up ps logs
.PHONY: all build run shell package publish_on_s3 info


shellmysql:
	docker run -it --rm orchardup/mysql mysql -h$(MYSQL_IP)


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
	$(PKGLVL) build


build ps:: before
	$(DOCKER_COMPOSE) $@


logs:: up
	$(DOCKER_COMPOSE) $@ $(SECTIONS)



publish_on_s3: $(PACKAGE)
	$(eval RAND := $(shell head -c 128 /dev/urandom | tr -dc A-Za-z0-9))
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
