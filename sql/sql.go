package sql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"moul.io/zapgorm"

	"pathwar.pw/entity"
)

type SQLOptions struct {
	Path string
}

type dump struct {
	Levels []*entity.Level
}

func FromOpts(opts *SQLOptions) (*gorm.DB, error) {
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
	if err := db.AutoMigrate(entity.All()...).Error; err != nil {
		return nil, err
	}
	// FIXME: apply real migrations

	return db, nil
}

func sqlDump(opts *SQLOptions) error {
	d := dump{}

	db, err := FromOpts(opts)
	if err != nil {
		return err
	}

	if err := db.Find(&d.Levels).Error; err != nil {
		return err
	}

	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
