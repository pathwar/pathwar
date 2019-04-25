package server

import "github.com/jinzhu/gorm"

type svc struct {
	jwtKey []byte
	db     *gorm.DB
}
