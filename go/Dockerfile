# builder
FROM            golang:1.13-alpine as builder
RUN             apk --no-cache --update add nodejs-npm make gcc g++ musl-dev openssl-dev git
ENV             GO111MODULE=on GOPROXY=https://proxy.golang.org,direct
COPY            go.mod go.sum /go/src/pathwar.land/go/
WORKDIR         /go/src/pathwar.land/go
RUN             go mod download
COPY            . .
RUN             wget -O ../rules.mk https://raw.githubusercontent.com/moul/rules.mk/master/rules.mk
RUN             make install

# runtime
FROM            alpine:3.10
RUN             apk --no-cache --update add openssl wget bash
RUN             wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
COPY            --from=builder /go/bin/pathwar /bin/pathwar
ENTRYPOINT      ["/bin/pathwar"]
CMD             ["api", "server"]
EXPOSE          8000 9111
