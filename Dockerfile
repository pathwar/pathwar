# builder
FROM            golang:1.17-alpine as builder
RUN             apk --no-cache --update add npm make gcc g++ musl-dev openssl-dev git perl-utils curl
WORKDIR         /go/src/pathwar.land
ENV             GO111MODULE=on GOPROXY=https://proxy.golang.org,direct
COPY            go.mod go.sum ./
RUN             go mod download
COPY            . .
WORKDIR         ./go
RUN             make install

# runtime
FROM            alpine:3.10
RUN             apk --no-cache --update add openssl wget bash
RUN             wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
COPY            --from=builder /go/bin/pathwar /bin/pathwar
ENTRYPOINT      ["/bin/pathwar"]
CMD             ["api", "server"]
EXPOSE          8000 9111
