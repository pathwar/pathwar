##
## functions
##

rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

##
## vars
##

GOPATH ?= $(HOME)/go
BIN = $(GOPATH)/bin/pathwar.pw
SOURCES = $(call rwildcard, ./, *.go)
PWCTL_SOURCES = $(call rwildcard,pwctl//,*.go)
OUR_SOURCES = $(filter-out $(call rwildcard,vendor//,*.go),$(SOURCES))
PROTOS = $(call rwildcard, ./, *.proto)
OUR_PROTOS = $(filter-out $(call rwildcard,vendor//,*.proto),$(PROTOS))
GENERATED_PB_FILES = \
	$(patsubst %.proto,%.pb.go,$(PROTOS)) \
	$(call rwildcard ./, *.gen.go)
PWCTL_OUT_FILES = \
	./pwctl/out/pwctl-linux-amd64
GENERATED_FILES = \
	$(GENERATED_PB_FILES) \
	$(PWCTL_OUT_FILES) \
	swagger.yaml
PROTOC_OPTS = -I/protobuf:vendor/github.com/grpc-ecosystem/grpc-gateway:vendor:.
RUN_OPTS ?=

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

.PHONY: run
run: $(BIN) mysql.up
	$(BIN) server $(RUN_OPTS)

.PHONY: install
install: $(BIN)
$(BIN): .proto.generated $(PWCTL_OUT_FILES) $(OUR_SOURCES)
	go install -v

.PHONY: mysql.up
mysql.up:
	docker-compose up -d mysql
	@echo "Waiting for mysql to be ready..."
	@while ! mysqladmin ping -h127.0.0.1 -P3306 --silent; do sleep 1; done
	@echo "Done."

.PHONY: mysql.down
mysql.down:
	docker-compose stop mysql || true
	docker-compose rm -f -v mysql || true

.PHONY: mysql.shell
mysql.shell:
	mysql -h127.0.0.1 -P3306 -uroot -puns3cur3 pathwar

.PHONY: mysql.dump
mysql.dump:
	mysqldump -h127.0.0.1 -P3306 -uroot -puns3cur3 pathwar

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
	go mod vendor
	docker run \
	  --user="$(shell id -u)" \
	  --volume="$(PWD):/go/src/pathwar.pw" \
	  --workdir="/go/src/pathwar.pw" \
	  --entrypoint="sh" \
	  --rm \
	  pathwar/protoc:v1 \
	  -xec "make _proto_generate"
	touch $@

.PHONY: _generate
_proto_generate: $(GENERATED_PB_FILES) swagger.yaml

$(PWCTL_OUT_FILES): $(PWCTL_SOURCES)
	mkdir -p ./pwctl/out
	GOOS=linux GOARCH=amd64 go build -o ./pwctl/out/pwctl-linux-amd64 ./pwctl/

.PHONY: test
test: .proto.generated
	go test -v ./...

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
	docker-compose exec server ./wait-for-it.sh mysql:3306 -- echo mysql ready
	docker-compose exec server ./wait-for-it.sh localhost:9111 -- echo gRPC ready
	sleep 5
	docker-compose exec server pathwar.pw sql adduser --sql-config=$$SQL_CONFIG --email=integration@example.com --username=integration --password=integration
	docker-compose run web npm test

.PHONY: lint
lint:
	golangci-lint run --verbose ./...

.PHONY: generate-fake-data
generate-fake-data:
	AUTH_TOKEN=`http --check-status :8000/authenticate username=integration | jq -r .token` && \
	  http POST :8000/dev/generate-fake-data Authorization:$$AUTH_TOKEN && \
	  http POST :8000/dev/sql-dump Authorization:$$AUTH_TOKEN
