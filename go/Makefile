## moul/rules.mk
GOPATH ?= $(HOME)/go
GOMOD_DIRS = ..
GOPKG =	pathwar.land
GOBINS = ./cmd/pathwar ./cmd/pwinit
DOCKER_IMAGE = pathwar/pathwar
DOCKERFILE_PATH = ../Dockerfile

GO_INSTALL_OPTS = -v -ldflags "-s -w -X pathwar.land/pathwar/v2/go/pkg/pwversion.Version=`git describe --tags --abbrev` -X pathwar.land/go/pkg/pwversion.Commit=`git rev-parse HEAD` -X pathwar.land/go/pkg/pwversion.Date=`date +%s` -X pathwar.land/go/pkg/pwversion.BuiltBy=makefile"
#PRE_INSTALL_STEPS += generate
#PRE_UNITTEST_STEPS += generate
#PRE_TEST_STEPS += generate
PRE_BUILD_STEPS += generate
PRE_LINT_STEPsS += generate
PRE_TIDY_STEPS += generate
PRE_BUMPDEPS_STEPS += generate
PRE_INSTALL_STEPS += ./out/pwinit-linux-amd64 ./out/docker-compose-dab

DEV_BIND ?= :8001
AGENT_NGINX_PORT ?= 8002
COMMON_RUN_OPTS ?= --debug --zipkin-endpoint=http://localhost:9411/api/v2/spans
API_DEV_OPTS ?=
AGENT_DEV_OPTS ?=

COMPILEDAEMON_OPTIONS ?= -exclude-dir=.git -color=true -build=go\ install -build-dir=./cmd/pathwar

COMPOSE_OPTS = -f docker-compose.yml -f docker-compose.dev.yml
include ../rules.mk  # see https://github.com/moul/rules.mk
##

help::
	@echo ""
	@echo "Pathwar commands:"
	@echo "  local-api          start an API (in auto-reload mode)"
	@echo "  local-agent        start an agent connected to the local API (in auto-reload mode)"
	@echo "  local-agent.auth   should be done once before local-agent"
	@echo "  up                 start everything in docker-compose"
	@echo ""

.PHONY: local-api
local-api: mysql.up zipkin.up generate
	$(GO) install github.com/githubnemo/CompileDaemon
	$(GO) install $(GO_INSTALL_OPTS) ./cmd/pathwar
	CompileDaemon $(COMPILEDAEMON_OPTIONS) -command="pathwar $(COMMON_RUN_OPTS) api --sso-unsafe server --with-pprof --bind=$(DEV_BIND) $(API_DEV_OPTS)"

.PHONY: api-down
api-down: mysql.down zipkin.down

.PHONY: local-agent
local-agent: install
	$(GO) install github.com/githubnemo/CompileDaemon
	CompileDaemon $(COMPILEDAEMON_OPTIONS) -command="pathwar $(COMMON_RUN_OPTS) agent --http-api-addr=http://localhost:8001 --port=$(AGENT_NGINX_PORT) $(AGENT_DEV_OPTS)"

.PHONY: local-agent.auth
local-agent.auth: install
	pathwar --debug agent --no-run

.PHONY: preprod-agent
preprod-agent: install
	pathwar $(COMMON_RUN_OPTS) agent --port=$(AGENT_NGINX_PORT) $(AGENT_DEV_OPTS)

.PHONY: agent-down
agent-down:
	pathwar --debug compose down --with-nginx

.PHONY: agent-nginx-config
agent-nginx-config:
	docker exec pathwar-agent-nginx cat /etc/nginx/nginx.conf

./out/docker-compose-dab:
	curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-`uname -s`-`uname -m`" -o $@
	chmod +x $@
	@#./$@ --version


.PHONY: ./out/pwinit-linux-amd64
./out/pwinit-linux-amd64:
	@mkdir -p out
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-linkmode external -extldflags -static" -mod=readonly -o $@ ./cmd/pwinit

.PHONY: up
up:
	docker-compose $(COMPOSE_OPTS) up -d

.PHONY: mysql.up
mysql.up:
	docker-compose $(COMPOSE_OPTS) up -d mysql

.PHONY: mysql.flush
mysql.flush: mysql.down
	docker volume rm -f go_pw_mysql_data

.PHONY: mysql.down
mysql.down:
	docker-compose stop mysql || true
	docker-compose rm -f -v mysql || true

.PHONY: mysql.logs
mysql.logs:
	docker-compose logs --tail=1000 -f mysql

.PHONY: mysql.shell
mysql.shell:
	mysql $(MYSQL_CONFIG) pathwar

.PHONY: mysql.dump
mysql.dump:
	mysqldump $(MYSQL_CONFIG) pathwar

.PHONY: zipkin.up
zipkin.up:
	docker-compose $(COMPOSE_OPTS) up -d zipkin

.PHONY: zipkin.logs
zipkin.logs:
	docker-compose $(COMPOSE_OPTS) logs --tail=1000 -f zipkin

.PHONY: zipkin.down
zipkin.down:
	docker-compose $(COMPOSE_OPTS) stop zipkin || true
	docker-compose $(COMPOSE_OPTS) rm -f -v zipkin || true

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

.PHONY: clean
clean:
	rm -f gen.sum $(wildcard */*/*.pb.go */*/*.pb.gw.go) $(wildcard out/*) $(wildcard */*/packrd/*) $(wildcard */*/*-packr.go)

.PHONY: packr
packr: ./out/pwinit-linux-amd64
	cd pkg/pwinit && packr2

.PHONY: install-light
install-light:
	go install ./cmd/pathwar
