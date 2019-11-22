module pathwar.land/go

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gobuffalo/packr/v2 v2.6.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.11.3
	github.com/jinzhu/gorm v1.9.11
	github.com/keycloak/kcinit v0.0.0-20181010192927-f85c3c5390ea
	github.com/oklog/run v1.0.0
	github.com/olekukonko/tablewriter v0.0.3
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/peterbourgon/ff v1.6.1-0.20190916204019-6cd704ec2eeb
	github.com/rs/cors v1.7.0
	github.com/treastech/logger v0.0.0-20180705232552-e381e9ecf2e3
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20190930134127-c5a3c61f89f3 // indirect
	golang.org/x/sys v0.0.0-20190927073244-c990c680b611 // indirect
	google.golang.org/appengine v1.6.4 // indirect
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.24.0
	gopkg.in/gormigrate.v1 v1.6.0
	gopkg.in/yaml.v3 v3.0.0-20191010095647-fc94e3f71652
	moul.io/godev v1.3.0
	moul.io/zapgorm v0.0.0-20190706070406-8138918b527b
)

replace (
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	gopkg.in/jcmturner/rpc.v1 => gopkg.in/jcmturner/rpc.v1 v1.1.0
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.13
