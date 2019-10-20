## moul/rules.mk
GOPATH ?= $(HOME)/go
GOPKG =		pathwar.land/go
GOBINS =	./cmd/pathwar ./cmd/pwchallenge
DOCKER_IMAGE =	pathwar/pathwar

GO_INSTALL_OPTS = -v -ldflags "-s -w -X pathwar.land/go/pkg/pwversion.Version=`git describe --tags --abbrev` -X pathwar.land/go/pkg/pwversion.Commit=`git rev-parse HEAD` -X pathwar.land/go/pkg/pwversion.Date=`date +%s` -X pathwar.land/go/pkg/pwversion.BuiltBy=makefile"
PRE_INSTALL_STEPS += generate
PRE_TEST_STEPS += generate
PRE_BUILD_STEPS += generate
PRE_LINT_STEPS += generate
PRE_TIDY_STEPS += generate
PRE_BUMPDEPS_STEPS += generate
PRE_INSTALL_STEPS += ./out/pwchallenge-linux-amd64
include ../rules.mk  # see https://github.com/moul/rules.mk
##

.PHONY: ./out/pwchallenge-linux-amd64
./out/pwchallenge-linux-amd64:
	mkdir -p out
	GOOS=linux GOARCH=amd64 $(GO) build -mod=readonly -o $@ ./cmd/pwchallenge

##
## mysql
##

MYSQL_CONFIG ?=	-h127.0.0.1 -P3306 -uroot -puns3cur3

.PHONY: run
run: install mysql.up
	pathwar --debug engine --sso-unsafe server --with-pprof

.PHONY: mysql.up
mysql.up:
	docker-compose up -d mysql
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
	docker-compose logs -f mysql

.PHONY: mysql.shell
mysql.shell:
	mysql $(MYSQL_CONFIG) pathwar

.PHONY: mysql.dump
mysql.dump:
	mysqldump $(MYSQL_CONFIG) pathwar

##
## generate
##

PROTOS_SRC := $(wildcard ../api/*.proto) $(wildcard ../api/internal/*.proto)
GEN_SRC := $(PROTOS_SRC) Makefile
.PHONY: generate
generate: gen.sum
gen.sum: $(GEN_SRC)
	shasum $(GEN_SRC) | sort > gen.sum.tmp
	diff -q gen.sum gen.sum.tmp || ( \
	  set -e; \
	  GO111MODULE=on go mod vendor; \
	  docker run \
	    --user=`id -u` \
	    --volume="$(PWD)/..:/go/src/pathwar.land" \
	    --workdir="/go/src/pathwar.land/go" \
	    --entrypoint="sh" \
	    --rm \
	    pathwar/protoc:4 \
	    -xec 'make generate_local'; \
	    make tidy \
	)

.PHONY: generate_local
generate_local:
	@set -e; for proto in $(PROTOS_SRC); do ( set -xe; \
	  protoc -I ../go/vendor/github.com/grpc-ecosystem/grpc-gateway:../api:../go/vendor:/protobuf --grpc-gateway_out=logtostderr=true:"$(GOPATH)/src" --gogofaster_out="plugins=grpc:$(GOPATH)/src" "$$proto" \
	); done
	goimports -w ./pkg ./cmd ./internal
	shasum $(GEN_SRC) | sort > gen.sum.tmp
	mv gen.sum.tmp gen.sum

.PHONY: clean
clean:
	rm -f gen.sum $(wildcard */*/*.pb.go */*/*.pb.gw.go) $(wildcard out/*)
