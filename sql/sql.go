package sql

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // required by gorm
	"go.uber.org/zap"
	"moul.io/zapgorm"

	"pathwar.pw/entity"
)

func FromOpts(opts *Options) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", opts.Path)
	if err != nil {
		return nil, err
	}

	log.SetOutput(ioutil.Discard)
	db.Callback().Create().Remove("gorm:update_time_stamp")
	db.Callback().Update().Remove("gorm:update_time_stamp")
	log.SetOutput(os.Stderr)

	db.SetLogger(zapgorm.New(zap.L().Named("vendor.gorm")))
	db = db.Set("gorm:auto_preload", true)
	db = db.Set("gorm:association_autoupdate", true)
	db.BlockGlobalUpdate(true)
	db.SingularTable(true)
	db.LogMode(true)
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(entity.All()...).Error; err != nil {
		return nil, err
	}
	// FIXME: use gormigrate

	return db, nil
}
