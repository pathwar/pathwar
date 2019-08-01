package server

import (
	"time"

	"github.com/jinzhu/gorm"
)
import "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"

var _ = options.E_Openapiv2Swagger

type svc struct {
	jwtKey    []byte
	db        *gorm.DB
	startedAt time.Time
}
