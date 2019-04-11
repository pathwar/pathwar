package server

import (
	"context"
	"math/rand"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/moul/gofakeit"
	"github.com/pkg/errors"

	"pathwar.pw/entity"
	"pathwar.pw/sql"
)

type svc struct {
	jwtKey []byte
	db     *gorm.DB
}

func (s *svc) Ping(_ context.Context, _ *Void) (*Void, error) {
	return &Void{}, nil
}

func (s *svc) Authenticate(ctx context.Context, input *AuthenticateInput) (*AuthenticateOutput, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// FIXME: use mapstructure
		"username": input.Username,
		// FIXME: if needed encrypt sensitive data
	})
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}
	// FIXME: set "cookie"
	return &AuthenticateOutput{
		Token: tokenString,
	}, nil
}

func (s *svc) UserSession(ctx context.Context, _ *Void) (*entity.UserSession, error) {
	sess, err := userSessionFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get context session")
	}
	return &sess, nil
}

func (s *svc) GenerateFakeData(ctx context.Context, _ *Void) (*Void, error) {
	hypervisors := []*entity.Hypervisor{}
	for i := 0; i < 3; i++ {
		hypervisor := &entity.Hypervisor{
			Name:    gofakeit.HipsterWord(),
			Address: gofakeit.IPv4Address(),
			Status:  entity.Hypervisor_Active,
		}
		hypervisors = append(hypervisors, hypervisor)
	}

	levels := []*entity.Level{}
	for i := 0; i < 5; i++ {
		level := &entity.Level{
			Name:        gofakeit.HipsterWord(),
			Description: gofakeit.HipsterSentence(10),
			Author:      gofakeit.Name(),
			Locale:      "fr_FR",
			IsDraft:     false,
			Versions:    []*entity.LevelVersion{},
		}
		for i := 0; i < 2; i++ {
			version := &entity.LevelVersion{
				Driver:    entity.LevelVersion_Docker,
				Version:   gofakeit.IPv4Address(),
				Changelog: gofakeit.HipsterSentence(5),
				IsDraft:   false,
				IsLatest:  i == 0,
				SourceURL: gofakeit.URL(),
				Flavors:   []*entity.LevelFlavor{},
			}
			for j := 0; j < 2; j++ {
				flavor := &entity.LevelFlavor{
					Instances: []*entity.LevelInstance{},
				}
				for k := 0; k < 2; k++ {
					instance := &entity.LevelInstance{
						HypervisorID: hypervisors[rand.Int()%len(hypervisors)].ID,
						Status:       entity.LevelInstance_Active,
					}
					flavor.Instances = append(flavor.Instances, instance)
				}
				version.Flavors = append(version.Flavors, flavor)
			}
			level.Versions = append(level.Versions, version)
		}
		levels = append(levels, level)
	}

	teams := []*entity.Team{}
	for i := 0; i < 5; i++ {
		team := &entity.Team{
			Name:        gofakeit.HipsterWord(),
			GravatarURL: gofakeit.URL(),
			Locale:      "fr_FR",
		}
		teams = append(teams, team)
	}

	users := []*entity.User{}
	for i := 0; i < 10; i++ {
		user := &entity.User{
			Username:    gofakeit.Name(),
			GravatarURL: gofakeit.URL(),
			WebsiteURL:  gofakeit.URL(),
			Locale:      "fr_FR",
			IsStaff:     false,
			Memberships: []*entity.TeamMember{},
		}
		users = append(users, user)
	}

	tournaments := []*entity.Tournament{}
	for i := 0; i < 3; i++ {
		tournament := &entity.Tournament{
			Name:       gofakeit.HipsterWord(),
			Status:     entity.Tournament_Started,
			Visibility: entity.Tournament_Public,
			IsDefault:  false,
		}
		tournaments = append(tournaments, tournament)
	}
	tournaments[0].IsDefault = true

	for _, entity := range hypervisors {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}
	for _, entity := range users {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}
	for _, entity := range levels {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}
	for _, entity := range tournaments {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}
	for _, entity := range teams {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}

	coupons := []*entity.Coupon{}
	for i := 0; i < 3; i++ {
		coupon := &entity.Coupon{
			Hash:               gofakeit.UUID(),
			MaxValidationCount: int32(rand.Int() % 5),
			Value:              int32(rand.Int() % 10),
			TournamentID:       tournaments[rand.Int()%len(tournaments)].ID,
		}
		coupons = append(coupons, coupon)
	}

	memberships := []*entity.TeamMember{}
	for _, user := range users {
		for i := 0; i < 2; i++ {
			memberships = append(
				memberships,
				&entity.TeamMember{
					TeamID: teams[rand.Int()%len(teams)].ID,
					UserID: user.ID,
					Role:   entity.TeamMember_Member, // or Owner
				},
			)
		}
	}

	for _, entity := range memberships {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}
	for _, entity := range coupons {
		if err := s.db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return nil, err
		}
	}
	return &Void{}, nil
}

func (s *svc) Dump(ctx context.Context, _ *Void) (*entity.Dump, error) {
	return sql.DoDump(s.db)
}
