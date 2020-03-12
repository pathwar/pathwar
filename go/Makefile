## moul/rules.mk
GOPATH ?= $(HOME)/go
GOMOD_DIR = ..
GOPKG =	pathwar.land
GOBINS = ./cmd/pathwar ./cmd/pwinit
DOCKER_IMAGE =pathwar/pathwar
DOCKERFILE_PATH = ../Dockerfile

GO_INSTALL_OPTS = -v -ldflags "-s -w -X pathwar.land/v2/go/pkg/pwversion.Version=`git describe --tags --abbrev` -X pathwar.land/go/pkg/pwversion.Commit=`git rev-parse HEAD` -X pathwar.land/go/pkg/pwversion.Date=`date +%s` -X pathwar.land/go/pkg/pwversion.BuiltBy=makefile"
PRE_INSTALL_STEPS += generate
PRE_UNITTEST_STEPS += generate
PRE_TEST_STEPS += generate
PRE_BUILD_STEPS += generate
PRE_LINT_STEPS += generate
PRE_TIDY_STEPS += generate
PRE_BUMPDEPS_STEPS += generate
PRE_INSTALL_STEPS += ./out/pwinit-linux-amd64

COMPOSE_OPTS = -f docker-compose.yml -f docker-compose.dev.yml
include ../rules.mk  # see https://github.com/moul/rules.mk
##

.PHONY: ./out/pwinit-linux-amd64
./out/pwinit-linux-amd64:
	mkdir -p out
	GOOS=linux GOARCH=amd64 $(GO) build -mod=readonly -o $@ ./cmd/pwinit

MYSQL_CONFIG ?=	-h127.0.0.1 -P3306 -uroot -puns3cur3

.PHONY: up
up:
	docker-compose $(COMPOSE_OPTS) up -d

.PHONY: run
run: install mysql.up
	pathwar --debug api --sso-unsafe server --with-pprof

.PHONY: mysql.up
mysql.up:
	docker-compose $(COMPOSE_OPTS) up -d mysql
	@echo "Waiting for mysql to be ready..."
	@while ! mysqladmin ping $(MYSQL_CONFIG) --silent; do sleep 1; done
	@echo "Done."

.PHONY: mysql.flush
mysql.flush: mysql.down
	docker volume rm -f go_mysql_data

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
	    pathwar/protoc:5 \
	    -xec 'make generate_local'; \
	    make tidy \
	)

.PHONY: generate_local
generate_local:
	@set -e; for proto in $(PROTOS_SRC); do ( set -xe; \
	  protoc -I ../vendor/github.com/grpc-ecosystem/grpc-gateway:../api:../vendor:/protobuf --grpc-gateway_out=logtostderr=true:"$(GOPATH)/src" --gogofaster_out="plugins=grpc:$(GOPATH)/src" "$$proto" \
	); done
	@### custom
	sed -i s@pathwar.land/go/pkg@pathwar.land/v2/go/pkg@ ./pkg/*/*.pb.go
	@### end of custom
	goimports -w ./pkg ./cmd ./internal
	shasum $(GEN_SRC) | sort > gen.sum.tmp
	mv gen.sum.tmp gen.sum

.PHONY: clean
clean:
	rm -f gen.sum $(wildcard */*/*.pb.go */*/*.pb.gw.go) $(wildcard out/*) $(wildcard */*/packrd/*) $(wildcard */*/*-packr.go)

.PHONY: packr
packr:
	cd pkg/pwinit && packr2
