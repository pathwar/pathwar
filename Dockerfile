# builder
FROM            golang:1.14-alpine as builder
RUN             apk --no-cache --update add nodejs-npm make gcc g++ musl-dev openssl-dev git perl-utils
RUN            	go get github.com/gobuffalo/packr/v2/packr2
ENV             GO111MODULE=on GOPROXY=https://proxy.golang.org,direct
COPY            go.mod go.sum /go/src/pathwar.land/
WORKDIR         /go/src/pathwar.land
RUN             go mod download
COPY            . .
WORKDIR         /go/src/pathwar.land/go
RUN             make packr
RUN             make install

# runtime
FROM            alpine:3.10
RUN             apk --no-cache --update add openssl wget bash
RUN             wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
COPY            --from=builder /go/bin/pathwar /bin/pathwar
ENTRYPOINT      ["/bin/pathwar"]
CMD             ["api", "server"]
EXPOSE          8000 9111
