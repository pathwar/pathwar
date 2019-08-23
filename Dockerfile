# builder
FROM            golang:1.12-alpine as builder
RUN             apk --no-cache --update add nodejs-npm make gcc g++ musl-dev openssl-dev git
ENV             GO111MODULE=on GOPROXY=https://goproxy.io
COPY            go.mod go.sum /go/src/pathwar.land/
WORKDIR         /go/src/pathwar.land
RUN             go mod download
COPY            . .
RUN             make _ci_prepare && make install

# runtime
FROM            alpine:3.8
RUN             apk --no-cache --update add openssl wget bash
RUN             wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
COPY            --from=builder /go/bin/pathwar.land /bin/pathwar.land
ENTRYPOINT      ["/bin/pathwar.land"]
CMD             ["server"]
EXPOSE          8000 9111
