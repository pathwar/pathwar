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
	$(PWCTL_OUT_FILES)
PROTOC_OPTS = -I/protobuf:vendor:.
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
run: $(BIN)
	$(BIN) server $(RUN_OPTS)

.PHONY: install
install: $(BIN)
$(BIN): .proto.generated $(PWCTL_OUT_FILES) $(OUR_SOURCES)
	go install -v

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
	rm -f $(GENERATED_FILES)
	go mod vendor
	docker run \
	  --user="$(shell id -u)" \
	  --volume="$(PWD):/go/src/pathwar.pw" \
	  --workdir="/go/src/pathwar.pw" \
	  --pwctl="sh" \
	  --rm \
	  pathwar/protoc:v1 \
	  -xec "make _proto_generate"
	touch $@

.PHONY: _generate
_proto_generate: $(GENERATED_PB_FILES)

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

.PHONY: docker.build
docker.build:
	docker build -t pathwar/pathwar .

.PHONY: docker.integration
docker.integration:
	docker-compose -f ./test/docker-compose.yml up -d server
	docker-compose -f ./test/docker-compose.yml run client
	docker-compose -f ./test/docker-compose.yml down

.PHONY: integration
integration:
	./test/integration.sh

.PHONY: lint
lint:
	golangci-lint run --verbose ./...
