module pathwar.land

require (
	cloud.google.com/go v0.44.3 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/GlenDC/go-external-ip v0.0.0-20170425150139-139229dcdddd
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/denisenkom/go-mssqldb v0.0.0-20190820223206-44cdfe8d8ba9 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20181024220401-bc4c1c238b55
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gobuffalo/logger v1.0.1 // indirect
	github.com/gobuffalo/packr/v2 v2.6.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190104160321-4832df01553a
	github.com/grpc-ecosystem/grpc-gateway v1.11.1
	github.com/jinzhu/gorm v1.9.10
	github.com/keycloak/kcinit v0.0.0-20181010192927-f85c3c5390ea
	github.com/lib/pq v1.2.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rogpeppe/go-internal v1.3.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190904154756-749cb33beabd // indirect
	google.golang.org/appengine v1.6.2 // indirect
	google.golang.org/genproto v0.0.0-20190905072037-92dd089d5514
	google.golang.org/grpc v1.23.0
	gotest.tools v0.0.0-20181223230014-1083505acf35 // indirect
	moul.io/zapgorm v0.0.0-20190706070406-8138918b527b
)

replace (
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	gopkg.in/jcmturner/rpc.v1 => gopkg.in/jcmturner/rpc.v1 v1.1.0
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.13
