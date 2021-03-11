package pwdb

import (
	"log"
	"reflect"
	"strings"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

// DefaultGormConfig that fits on both prod and test
var DefaultGormConfig = gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
}

func Configure(db *gorm.DB, sfn *snowflake.Node) (*gorm.DB, error) {
	err := db.Callback().Create().Before("gorm:create").Register("snowflake_id:before_create", snowFlakeIDS(sfn))
	if err != nil {
		return nil, errcode.ErrDBAddCallback.Wrap(err)
	}
	if err := migrate(db); err != nil {
		return nil, errcode.ErrDBRunMigrations.Wrap(err)
	}
	return db, nil
}

func snowFlakeIDS(sfn *snowflake.Node) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Schema != nil {
			for _, field := range db.Statement.Schema.PrimaryFields {
				// ensure that it's an ID
				if field.DataType == schema.Int && strings.Contains(field.Name, "ID") {
					switch db.Statement.ReflectValue.Kind() {
					case reflect.Slice, reflect.Array:
						for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
							if _, isZero := field.ValueOf(db.Statement.ReflectValue.Index(i)); isZero {
								err := field.Set(db.Statement.ReflectValue.Index(i), sfn.Generate().Int64())
								if err != nil {
									log.Println(errcode.ErrDBSetSnowflakeID.Wrap(err))
								}
							}
						}
					case reflect.Struct:
						if _, isZero := field.ValueOf(db.Statement.ReflectValue); isZero {
							err := field.Set(db.Statement.ReflectValue, sfn.Generate().Int64())
							if err != nil {
								log.Println(errcode.ErrDBSetSnowflakeID.Wrap(err))
							}
						}
					}
				}
			}
		}
	}
}
