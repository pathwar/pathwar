# builder
FROM            golang:1.11-alpine as builder
RUN             apk --no-cache --update add nodejs-npm make gcc g++ musl-dev openssl-dev git
ENV             GO111MODULE=on
COPY            go.* /go/src/pathwar.pw/
WORKDIR         /go/src/pathwar.pw
RUN             go get .
COPY            . .
RUN             touch .proto.generated && make install

# runtime
FROM            alpine:3.8
RUN             apk --no-cache --update add openssl
COPY            --from=builder /go/bin/pathwar.pw /bin/pathwar.pw
ENTRYPOINT      ["/bin/pathwar.pw"]
CMD             ["server"]
EXPOSE          8000 9111
