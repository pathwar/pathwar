package sql // import "pathwar.pw/sql"

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // required by gorm
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"moul.io/zapgorm"

	"pathwar.pw/entity"
)

func FromOpts(opts *Options) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", opts.Config)
	if err != nil {
		return nil, err
	}

	log.SetOutput(ioutil.Discard)
	//db.Callback().Create().Remove("gorm:update_time_stamp")
	//db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Callback().Create().Before("gorm:create").Register("pathwar_before_create", beforeCreate)
	log.SetOutput(os.Stderr)

	db.SetLogger(zapgorm.New(zap.L().Named("vendor.gorm")))
	db = db.Set("gorm:auto_preload", false)
	db = db.Set("gorm:association_autoupdate", false)
	db.BlockGlobalUpdate(true)
	db.SingularTable(true)
	db.LogMode(true)
	if err := db.AutoMigrate(entity.All()...).Error; err != nil {
		return nil, err
	}
	for _, fk := range entity.ForeignKeys() {
		e := entity.ByName(fk[0])
		if err := db.Model(e).AddForeignKey(fk[1], fk[2], "RESTRICT", "RESTRICT").Error; err != nil {
			return nil, err
		}
	}
	// FIXME: use gormigrate

	return db, nil
}

func beforeCreate(scope *gorm.Scope) {
	id, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		panic(err)
	}
	out := base64.StdEncoding.EncodeToString(id)
	if err := scope.SetColumn("ID", out); err != nil {
		panic(err)
	}
}
