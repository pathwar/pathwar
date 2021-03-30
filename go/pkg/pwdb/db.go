package pwdb

import (
	"log"
	"reflect"
	"strings"
	time "time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func Configure(db *gorm.DB, sfn *snowflake.Node) (*gorm.DB, error) {
	err := db.Callback().Delete().Replace("gorm:delete", Delete)
	if err != nil {
		return nil, errcode.ErrDBAddCallback.Wrap(err)
	}
	err = db.Callback().Query().Before("gorm:query").Register("softDelete:before_query", beforeQuery)
	if err != nil {
		return nil, errcode.ErrDBAddCallback.Wrap(err)
	}
	err = db.Callback().Create().Before("gorm:create").Register("snowflake_id:before_create", snowFlakeIDS(sfn))
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

// simulate softDelete query system
func beforeQuery(db *gorm.DB) {
	if db.Statement.Schema != nil && !db.Statement.Unscoped {
		deletedAtField := db.Statement.Schema.LookUpField("deleted_at")
		if deletedAtField != nil {
			db.Statement.AddClause(clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{
						Column: clause.Column{Table: clause.CurrentTable, Name: "deleted_at"},
						Value:  nil,
					}},
			})
		}
	}
}

// simulate softDelete delete system
// !! need manual update in case gorm official Delete callback changes !!
func Delete(db *gorm.DB) {
	if db.Error == nil {
		if db.Statement.Schema != nil && !db.Statement.Unscoped {
			for _, c := range db.Statement.Schema.DeleteClauses {
				db.Statement.AddClause(c)
			}
		}

		if db.Statement.SQL.String() == "" {
			db.Statement.SQL.Grow(100)
			db.Statement.AddClauseIfNotExists(clause.Delete{})

			if db.Statement.Schema != nil {
				_, queryValues := schema.GetIdentityFieldValuesMap(db.Statement.ReflectValue, db.Statement.Schema.PrimaryFields)
				column, values := schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}

				if db.Statement.ReflectValue.CanAddr() && db.Statement.Dest != db.Statement.Model && db.Statement.Model != nil {
					_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(db.Statement.Model), db.Statement.Schema.PrimaryFields)
					column, values = schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

					if len(values) > 0 {
						db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
					}
				}
			}

			db.Statement.AddClauseIfNotExists(clause.From{})
			// here is the hacky softDelete part
			deletedAtField := db.Statement.Schema.LookUpField("deleted_at")
			if deletedAtField != nil && !db.Statement.Unscoped {
				db.Statement.AddClauseIfNotExists(clause.Update{})
				db.Statement.AddClause(clause.Set([]clause.Assignment{
					{
						Column: clause.Column{Name: "deleted_at"},
						Value:  time.Now(),
					},
				}))
				db.Statement.Build("UPDATE", "SET", "WHERE")
			} else {
				db.Statement.Build("DELETE", "FROM", "WHERE")
			}
		}

		if _, ok := db.Statement.Clauses["WHERE"]; !db.AllowGlobalUpdate && !ok && db.Error == nil {
			_ = db.AddError(gorm.ErrMissingWhereClause)
			return
		}

		if !db.DryRun && db.Error == nil {
			result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)

			if err == nil {
				db.RowsAffected, _ = result.RowsAffected()
			} else {
				_ = db.AddError(err)
			}
		}
	}
}
