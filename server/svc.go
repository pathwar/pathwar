package server

import (
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	"github.com/jinzhu/gorm"
	"pathwar.land/client"
)

var _ = options.E_Openapiv2Swagger

type svc struct {
	db        *gorm.DB
	client    client.Options
	startedAt time.Time
}
