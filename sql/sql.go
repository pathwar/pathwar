package sql // import "pathwar.land/sql"

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
)

func FromOpts(opts *Options) (*gorm.DB, error) {
	sqlConfig := opts.Config
	if envConfig := os.Getenv("SQL_CONFIG"); envConfig != "" { // this should be done using viper's built-in env support
		sqlConfig = envConfig
	}
	zap.L().Debug("opening sql", zap.String("config", sqlConfig))
	db, err := gorm.Open("mysql", sqlConfig)
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

	if err := migrate(db); err != nil {
		return nil, err
	}

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
