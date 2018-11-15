FROM    moul/protoc-gen-gotemplate:latest as pgg

FROM    golang:1.11-alpine as builder
RUN     apk --no-cache add make git go rsync libc-dev openssh docker
RUN     go get -u -v \
        github.com/opsee/protobuf/protoc-gen-gogoopsee

FROM    golang:1.11-alpine
RUN     apk --no-cache add git make protobuf gcc libc-dev && \
        mkdir -p /.cache/go-build && \
        chmod -R 777 /.cache
COPY    --from=pgg /go/bin/* /go/bin/
COPY    --from=builder /go/bin/* /go/bin/
COPY    --from=pgg /protobuf /protobuf
ENV     GOPATH=/go \
        PATH=/go/bin:${PATH} \
        GOROOT=/usr/local/go
