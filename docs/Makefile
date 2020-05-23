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
	  (set -e; cd ../go; GO111MODULE=on go mod vendor;); \
	  docker run \
	    --user=`id -u` \
	    --volume="$(PWD)/..:/go/src/pathwar.land" \
	    --workdir="/go/src/pathwar.land/docs" \
	    --entrypoint="sh" \
	    --rm \
	    pathwar/protoc:6 \
	    -xec 'make generate_local'; \
	)

.PHONY: generate_local
generate_local:
	protoc -I ../vendor/github.com/grpc-ecosystem/grpc-gateway:../api:../vendor:/protobuf --swagger_out=logtostderr=true:. ../api/pwapi.proto
	echo 'swagger: "2.0"' > swagger.yaml.tmp
	cat ./pwapi.swagger.json | json2yaml | grep -v 'swagger: "2.0"' >> swagger.yaml.tmp
	rm -f ./pwapi.swagger.json
	mv swagger.yaml.tmp ../swagger.yaml
	eclint fix ../swagger.yaml
	shasum $(GEN_SRC) | sort > gen.sum.tmp
	mv gen.sum.tmp gen.sum

.PHONY: clean
clean:
	rm -f gen.sum swagger.yaml
