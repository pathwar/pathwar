## moul/rules.mk
GOPATH ?= $(HOME)/go
GOMOD_DIR = ..
GOPKG =	pathwar.land
GOBINS = ./cmd/pathwar ./cmd/pwinit
DOCKER_IMAGE = pathwar/pathwar
DOCKERFILE_PATH = ../Dockerfile
GO_INSTALL_OPTS = -v -ldflags "-s -w -X pathwar.land/pathwar/v2/go/pkg/pwversion.Version=`git describe --tags --abbrev` -X pathwar.land/go/pkg/pwversion.Commit=`git rev-parse HEAD` -X pathwar.land/go/pkg/pwversion.Date=`date +%s` -X pathwar.land/go/pkg/pwversion.BuiltBy=makefile"
PRE_INSTALL_STEPS += generate
PRE_UNITTEST_STEPS += generate
PRE_TEST_STEPS += generate
PRE_BUILD_STEPS += generate
PRE_LINT_STEPsS += generate
PRE_TIDY_STEPS += generate
PRE_BUMPDEPS_STEPS += generate
PRE_INSTALL_STEPS += ./out/pwinit-linux-amd64
include ../rules.mk  # see https://github.com/moul/rules.mk
## EO-moul/rules.mk

## Dev commands
COMPILEDAEMON_OPTIONS ?= -exclude-dir=.git -color=true -build=go\ install -build-dir=./cmd/pathwar
AGENT_DEV_OPTS ?=
COMMON_RUN_OPTS ?= --debug
API_DEV_OPTS ?=
DEV_BIND ?= :8001
AGENT_NGINX_PORT ?= 8002
api.run: generate
	$(GO) install github.com/githubnemo/CompileDaemon
	$(GO) install $(GO_INSTALL_OPTS) ./cmd/pathwar
	CompileDaemon $(COMPILEDAEMON_OPTIONS) -command="pathwar $(COMMON_RUN_OPTS) api --sso-unsafe server --with-pprof --bind=$(DEV_BIND) $(API_DEV_OPTS)"
agent.run: install
	$(GO) install github.com/githubnemo/CompileDaemon
	CompileDaemon $(COMPILEDAEMON_OPTIONS) -command="pathwar $(COMMON_RUN_OPTS) agent --http-api-addr=http://localhost:8001 --port=$(AGENT_NGINX_PORT) $(AGENT_DEV_OPTS)"
agent.auth: install
	pathwar --debug agent --no-run
agent.preprod: install
	pathwar $(COMMON_RUN_OPTS) agent --port=$(AGENT_NGINX_PORT) $(AGENT_DEV_OPTS)
agent.down:
	pathwar --debug compose down --with-nginx
agent.nginx-config:
	docker exec pathwar-agent-nginx cat /etc/nginx/nginx.conf

##
## generate
##
PROTOS_SRC := $(wildcard ../api/*.proto) $(wildcard ../api/internal/*.proto)
GEN_SRC := $(PROTOS_SRC) Makefile
.PHONY: generate
generate: gen.sum
gen.sum: $(GEN_SRC)
	@shasum $(GEN_SRC) | sort > gen.sum.tmp
	@diff -q gen.sum gen.sum.tmp || ( \
	  set -xe; \
	  GO111MODULE=on go mod vendor; \
	  docker run \
	    --user=`id -u` \
	    --volume="$(PWD)/..:/go/src/pathwar.land" \
	    --workdir="/go/src/pathwar.land/go" \
	    --entrypoint="sh" \
	    --rm \
	    pathwar/protoc:6 \
	    -xec 'make generate_local'; \
	    make tidy \
	)
.PHONY: generate_local
generate_local:
	@set -e; for proto in $(PROTOS_SRC); do ( set -xe; \
	  protoc -I ../vendor/github.com/grpc-ecosystem/grpc-gateway:../api:../vendor:/protobuf --grpc-gateway_out=logtostderr=true:"$(GOPATH)/src" --gogofaster_out="plugins=grpc:$(GOPATH)/src" "$$proto" \
	); done
	@### custom
	sed -i s@pathwar.land/go/pkg@pathwar.land/pathwar/v2/go/pkg@ ./pkg/*/*.pb.go
	@### end of custom
	goimports -w ./pkg ./cmd ./internal
	shasum $(GEN_SRC) | sort > gen.sum.tmp
	mv gen.sum.tmp gen.sum

## Helpers
.PHONY: clean
clean:
	rm -f gen.sum $(wildcard */*/*.pb.go */*/*.pb.gw.go) $(wildcard out/*) $(wildcard */*/packrd/*) $(wildcard */*/*-packr.go)

.PHONY: _docker
_docker: ./out/pwinit-linux-amd64
	cd pkg/pwinit && packr2

.PHONY: ./out/pwinit-linux-amd64
./out/pwinit-linux-amd64:
	mkdir -p out
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-linkmode external -extldflags -static" -mod=readonly -o $@ ./cmd/pwinit
