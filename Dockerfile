# builder
FROM            golang:1.12-alpine as builder
RUN             apk --no-cache --update add nodejs-npm make gcc g++ musl-dev openssl-dev git
ENV             GO111MODULE=on GOPROXY=https://goproxy.io
COPY            go.mod go.sum /go/src/pathwar.pw/
WORKDIR         /go/src/pathwar.pw
RUN             go mod download
COPY            . .
RUN             touch .proto.generated && make install

# runtime
FROM            alpine:3.8
RUN             apk --no-cache --update add openssl wget bash
RUN             wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
COPY            --from=builder /go/bin/pathwar.pw /bin/pathwar.pw
ENTRYPOINT      ["/bin/pathwar.pw"]
CMD             ["server"]
EXPOSE          8000 9111
