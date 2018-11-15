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
PROTOS = $(call rwildcard, ./, *.proto)
GENERATED_FILES = \
	$(patsubst %.proto,%.pb.go,$(PROTOS)) \
	$(call rwildcard ./, *.gen.go)
PROTOC_OPTS = -I/protobuf:.
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
$(BIN): .generated $(SOURCES)
	go install -v

.PHONY: clean
clean:
	rm -f $(GENERATED_FILES) .generated

.PHONY: generate
generate: .generated
.generated: $(PROTOS)
	rm -f $(GENERATED_FILES)
	docker run \
	  --user="$(shell id -u)" \
	  --volume="$(PWD):/go/src/pathwar.pw" \
	  --workdir="/go/src/pathwar.pw" \
	  --entrypoint="sh" \
	  --rm \
	  moul/protoc-gen-gotemplate \
	  -xec "make _generate"
	touch $@

.PHONY: _generate
_generate: $(GENERATED_FILES)

%.pb.go: %.proto
	protoc $(PROTOC_OPTS) --gofast_out=plugins=grpc:"$(GOPATH)/src" "$(dir $<)"/*.proto
