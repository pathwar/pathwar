package sql

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.land/entity"
	"pathwar.land/pkg/cli"
)

type dumpOptions struct {
	sql Options `mapstructure:"sql"`

	// additional dump filters
	// --anonymize
}

type dumpCommand struct{ opts dumpOptions }

func (cmd *dumpCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "dump",
		RunE: func(_ *cobra.Command, args []string) error {
			opts := cmd.opts
			opts.sql = GetOptions(commands)
			return runDump(&opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	commands["sql"].ParseFlags(cc.Flags())
	return cc
}
func (cmd *dumpCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *dumpCommand) ParseFlags(flags *pflag.FlagSet) {
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func DoDump(db *gorm.DB) (*entity.Dump, error) {
	dump := entity.Dump{}
	if err := db.Find(&dump.Levels).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.LevelVersions).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.LevelFlavors).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.LevelInstances).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Hypervisors).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.UserSessions).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Users).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Teams).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TeamMembers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Tournaments).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TournamentTeams).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TournamentMembers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Coupons).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Events).Error; err != nil {
		return nil, err
	}
	return &dump, nil
}

func runDump(opts *dumpOptions) error {
	db, err := FromOpts(&opts.sql)
	if err != nil {
		return err
	}

	dump, err := DoDump(db)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(dump, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
