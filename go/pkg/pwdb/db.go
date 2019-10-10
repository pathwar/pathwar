package pwdb

import (
	"encoding/base64"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"moul.io/zapgorm"
)

type Opts struct {
	Logger *zap.Logger

	// internal
	skipFK bool
}

func Configure(db *gorm.DB, opts Opts) (*gorm.DB, error) {
	db.SetLogger(zapgorm.New(opts.Logger))
	//db.Callback().Create().Remove("gorm:update_time_stamp")
	//db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Callback().Create().Before("gorm:create").Register("pathwar_before_create", beforeCreate)
	db = db.Set("gorm:auto_preload", false)
	db = db.Set("gorm:association_autoupdate", false)
	db.BlockGlobalUpdate(true)
	db.SingularTable(true)
	db.LogMode(true)
	if err := migrate(db, opts); err != nil {
		return nil, err
	}
	return db, nil
}

func beforeCreate(scope *gorm.Scope) {
	switch scope.TableName() {
	case "user":
		return
	}
	id, err := uuid.NewV4().MarshalBinary()
	if err != nil {
		panic(err)
	}
	out := base64.StdEncoding.EncodeToString(id)
	if err := scope.SetColumn("ID", out); err != nil {
		panic(err)
	}
}
