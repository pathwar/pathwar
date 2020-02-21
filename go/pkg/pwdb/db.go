package pwdb

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"moul.io/zapgorm"
	"pathwar.land/go/v2/pkg/errcode"
)

type Opts struct {
	Logger *zap.Logger

	// internal
	skipFK bool
}

func Configure(db *gorm.DB, sfn *snowflake.Node, opts Opts) (*gorm.DB, error) {
	db.SetLogger(zapgorm.New(opts.Logger))
	//db.Callback().Create().Remove("gorm:update_time_stamp")
	//db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Callback().Create().Before("gorm:create").Register("pathwar_before_create", beforeCreate(sfn))
	db = db.Set("gorm:auto_preload", false)
	db = db.Set("gorm:association_autoupdate", false)
	db.BlockGlobalUpdate(true)
	db.SingularTable(true)
	db.LogMode(true)
	if err := migrate(db, sfn, opts); err != nil {
		return nil, errcode.ErrDBRunMigrations.Wrap(err)
	}
	return db, nil
}

func beforeCreate(sfn *snowflake.Node) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		id := sfn.Generate().Int64()
		if err := scope.SetColumn("ID", id); err != nil {
			panic(err)
		}
	}
}
