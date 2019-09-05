##
## functions
##

rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

##
## vars
##

GOPATH ?= $(HOME)/go
GO ?= go
BIN := $(GOPATH)/bin/pathwar.land
SOURCES := $(call rwildcard, ./, *.go)
PWCTL_SOURCES := $(call rwildcard,pwctl//,*.go)
OUR_SOURCES := $(filter-out $(call rwildcard,vendor//,*.go),$(SOURCES))
PROTOS := $(call rwildcard, ./, *.proto)
OUR_PROTOS := $(filter-out $(call rwildcard,vendor//,*.proto),$(PROTOS))
GENERATED_PB_FILES := \
	$(patsubst %.proto,%.pb.go,$(PROTOS)) \
	$(call rwildcard ./, *.gen.go)
PWCTL_OUT_FILES := \
	./pwctl/out/pwctl-linux-amd64
GENERATED_FILES := \
	$(GENERATED_PB_FILES) \
	$(PWCTL_OUT_FILES) \
	swagger.yaml
PROTOC_OPTS := -I/protobuf:vendor/github.com/grpc-ecosystem/grpc-gateway:vendor:.
RUN_OPTS ?=
SERVERDB_CONFIG ?=	-h127.0.0.1 -P3306 -uroot -puns3cur3

##
## rules
##

.PHONY: help
help:
	@echo "Make commands:"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: \
	  '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | \
	  sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | grep -v / | \
	  sed 's/^/  $(HELP_MSG_PREFIX)make /'

.PHONY: test
test: unittest lint tidy

.PHONY: run
run: $(BIN) serverdb.up
	$(BIN) server $(RUN_OPTS)

.PHONY: install
install: $(BIN)
$(BIN): .proto.generated $(PWCTL_OUT_FILES) $(OUR_SOURCES)
	$(GO) install -v -ldflags "-s -w -X pathwar.land/version.Version=`git describe --tags --abbrev` -X pathwar.land/version.Commit=`git rev-parse HEAD` -X pathwar.land/version.Date=`date +%s` -X pathwar.land/version.BuiltBy=makefile"


.PHONY: serverdb.up
serverdb.up:
	docker-compose up -d serverdb
	@echo "Waiting for serverdb to be ready..."
	@while ! mysqladmin ping $(SERVERDB_CONFIG) --silent; do sleep 1; done
	@echo "Done."

.PHONY: serverdb.flush
serverdb.flush: serverdb.down
	docker volume rm pathwarland_serverdb_data

.PHONY: serverdb.down
serverdb.down:
	docker-compose stop serverdb || true
	docker-compose rm -f -v serverdb || true

.PHONY: serverdb.logs
serverdb.logs:
	docker-compose logs -f serverdb

.PHONY: serverdb.shell
serverdb.shell:
	mysql $(SERVERDB_CONFIG) pathwar

.PHONY: serverdb.dump
serverdb.dump:
	mysqldump $(SERVERDB_CONFIG) pathwar

.PHONY: keycloakdb.shell
keycloakdb.shell:
	mysql -h127.0.0.1 -P3307 -uroot -puns3cur3 keycloak

.PHONY: keycloakdb.dump
keycloakdb.dump:
	mysqldump -h127.0.0.1 -P3307 -uroot -puns3cur3 keycloak

.PHONY: clean
clean:
	rm -f $(GENERATED_FILES) .proto.generated

.PHONY: _ci_prepare
_ci_prepare:
	mkdir -p $(dir $(GENERATED_FILES))
	touch $(OUR_PROTOS) $(GENERATED_FILES)
	sleep 1
	touch .proto.generated

.PHONY: generate
generate: .proto.generated
.proto.generated: $(OUR_PROTOS)
	rm -f $(GENERATED_PB_FILES)
	$(GO) mod vendor
	docker run \
	  --user="$(shell id -u)" \
	  --volume="$(PWD):/go/src/pathwar.land" \
	  --workdir="/go/src/pathwar.land" \
	  --entrypoint="sh" \
	  --rm \
	  pathwar/protoc:v2 \
	  -xec "make _proto_generate"
	touch $@
	rm -rf vendor

.PHONY: _generate
_proto_generate: $(GENERATED_PB_FILES) swagger.yaml

$(PWCTL_OUT_FILES): $(PWCTL_SOURCES)
	mkdir -p ./pwctl/out
	GOOS=linux GOARCH=amd64 $(GO) build -mod=readonly -o ./pwctl/out/pwctl-linux-amd64 ./pwctl/

.PHONY: unittest
unittest: .proto.generated
	echo "" > /tmp/coverage.txt
	set -e; for dir in `find . -type f -name "go.mod" | sed -r 's@/[^/]+$$@@' | sort | uniq`; do (set -xe; \
	  cd $$dir; \
	  $(GO) test -v -mod=readonly -cover -coverprofile=/tmp/profile.out -covermode=atomic -race ./...; \
	  if [ -f /tmp/profile.out ]; then \
	    cat /tmp/profile.out >> /tmp/coverage.txt; \
	    rm -f /tmp/profile.out; \
	  fi); done
	mv /tmp/coverage.txt .

%.pb.go: %.proto
	protoc \
	  $(PROTOC_OPTS) \
	  --grpc-gateway_out=logtostderr=true:"$(GOPATH)/src" \
	  --gogofaster_out=plugins=grpc:"$(GOPATH)/src" \
	  "$(dir $<)"/*.proto

swagger.yaml: $(PROTOS)
	protoc \
	  $(PROTOC_OPTS) \
	  --swagger_out=logtostderr=true:. \
	  ./server/*.proto
	echo 'swagger: "2.0"' > swagger.yaml.tmp
	cat server/server.swagger.json | json2yaml | grep -v 'swagger:."2.0"' >> swagger.yaml.tmp
	rm -f server/server.swagger.json
	mv swagger.yaml.tmp swagger.yaml
	eclint fix swagger.yaml

.PHONY: docker.build
docker.build:
	docker build -t pathwar/pathwar:latest .

.PHONY: integration
integration: integration.build integration.run

.PHONY: integration.build
integration.build:
	docker-compose build server web

.PHONY:integration.run
integration.run:
	docker-compose up -d --no-build server
	docker-compose exec server ./wait-for-it.sh serverdb:3306 -- echo serverdb ready
	docker-compose exec server ./wait-for-it.sh localhost:9111 -- echo gRPC ready
	sleep 5
	docker-compose exec server pathwar.land sql adduser --sql-config=$$SQL_CONFIG --email=integration@example.com --username=integration --password=integration
	docker-compose run web npm test

.PHONY: lint
lint:
	set -e; for dir in `find . -type f -name "go.mod" | sed 's@/[^/]*$$@@' | sort | uniq`; do (set -xe; \
	  cd $$dir; \
	  golangci-lint run --verbose ./...; \
	); done

.PHONY: tidy
tidy:
	set -e; for dir in `find . -type f -name "go.mod" | sed 's@/[^/]*$$@@' | sort | uniq`; do (set -xe; \
	  cd $$dir; \
	  $(GO) mod tidy; \
	); done

.PHONY: bump-go-deps
bump-go-deps:
	set -e; for dir in `find . -type f -name "go.mod" | sed 's@/[^/]*$$@@' | sort | uniq`; do (set -xe; \
	  cd $$dir; \
	  $(GO) get -u ./...; \
	); done

.PHONY: generate-fake-data
generate-fake-data:
	AUTH_TOKEN=`http --check-status :8000/authenticate username=integration | jq -r .token` && \
	  http POST :8000/dev/generate-fake-data Authorization:$$AUTH_TOKEN && \
	  http POST :8000/dev/sql-dump Authorization:$$AUTH_TOKEN
