package sql

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	//_ "github.com/mattn/go-sqlite3" // required by gorm
	_ "github.com/go-sql-driver/mysql" // required by gorm
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
	for _, fk := range [][3]string{
		{"Achievement", "team_member_id", "team_member(id)"},
		{"Coupon", "team_member_id", "team_member(id)"},
		{"LevelFlavor", "level_id", "level(id)"},
		{"LevelInstance", "hypervisor_id", "hypervisor(id)"},
		{"LevelInstance", "level_flavor_id", "level_flavor(id)"},
		{"LevelSubscription", "level_flavor_id", "level_flavor(id)"},
		{"LevelSubscription", "tournament_team_id", "tournament_team(id)"},
		{"Notification", "user_id", "user(id)"},
		{"ShopItem", "tournament_team_id", "tournament_team(id)"},
		{"TeamMember", "tournament_team_id", "tournament_team(id)"},
		{"TeamMember", "user_id", "user(id)"},
		{"TournamentTeam", "team_id", "team(id)"},
		{"TournamentTeam", "tournament_id", "tournament(id)"},
		{"UserSession", "user_id", "user(id)"},
		{"WhoswhoAttempt", "author_team_member_id", "team_member(id)"},
		{"WhoswhoAttempt", "target_team_member_id", "team_member(id)"},
		{"WhoswhoAttempt", "target_tournament_team_id", "tournament_team(id)"},
	} {
		e := entity.ByName(fk[0])
		if err := db.Model(e).AddForeignKey(fk[1], fk[2], "RESTRICT", "RESTRICT").Error; err != nil {
			return nil, err
		}
	}
	// FIXME: use gormigrate

	return db, nil
}
