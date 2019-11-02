package pwdb

func ByName(name string) interface{} {
	return AllMap()[name]
}

func AllMap() map[string]interface{} {
	return map[string]interface{}{
		"Achievement":           Achievement{},
		"Challenge":             Challenge{},
		"ChallengeFlavor":       ChallengeFlavor{},
		"ChallengeInstance":     ChallengeInstance{},
		"ChallengeSubscription": ChallengeSubscription{},
		"ChallengeValidation":   ChallengeValidation{},
		"Coupon":                Coupon{},
		"CouponValidation":      CouponValidation{},
		"Hypervisor":            Hypervisor{},
		"InventoryItem":         InventoryItem{},
		"Notification":          Notification{},
		"Organization":          Organization{},
		"OrganizationMember":    OrganizationMember{},
		"Season":                Season{},
		"SeasonChallenge":       SeasonChallenge{},
		"Team":                  Team{},
		"TeamMember":            TeamMember{},
		"User":                  User{},
		"WhoswhoAttempt":        WhoswhoAttempt{},
	}
}

func All() []interface{} {
	out := []interface{}{}
	for _, entry := range AllMap() {
		out = append(out, entry)
	}
	return out
}

func ForeignKeys() [][3]string {
	return [][3]string{
		// {"User", "active_team_member_id", "team_member(id)"}, // FIXME: check why this cause an error!?
		{"Achievement", "user_id", "user(id)"},
		{"Achievement", "team_id", "team(id)"},
		{"Achievement", "challenge_validation_id", "challenge_validation(id)"},
		{"ChallengeFlavor", "challenge_id", "challenge(id)"},
		{"ChallengeInstance", "flavor_id", "challenge_flavor(id)"},
		{"ChallengeInstance", "hypervisor_id", "hypervisor(id)"},
		{"ChallengeSubscription", "season_challenge_id", "season_challenge(id)"},
		{"ChallengeSubscription", "team_id", "team(id)"},
		{"ChallengeSubscription", "author_id", "user(id)"},
		{"ChallengeValidation", "challenge_subscription_id", "challenge_subscription(id)"},
		{"ChallengeValidation", "author_id", "user(id)"},
		{"Coupon", "season_id", "season(id)"},
		{"CouponValidation", "author_id", "user(id)"},
		{"CouponValidation", "team_id", "team(id)"},
		{"CouponValidation", "coupon_id", "coupon(id)"},
		{"InventoryItem", "owner_id", "team_member(id)"},
		{"Notification", "user_id", "user(id)"},
		{"OrganizationMember", "organization_id", "organization(id)"},
		{"OrganizationMember", "user_id", "user(id)"},
		{"SeasonChallenge", "flavor_id", "challenge_flavor(id)"},
		{"SeasonChallenge", "season_id", "season(id)"},
		{"Team", "organization_id", "organization(id)"},
		{"Team", "season_id", "season(id)"},
		{"TeamMember", "team_id", "team(id)"},
		{"TeamMember", "user_id", "user(id)"},
		{"WhoswhoAttempt", "author_id", "user(id)"},
		{"WhoswhoAttempt", "author_team_id", "team(id)"},
		{"WhoswhoAttempt", "target_user_id", "user(id)"},
		{"WhoswhoAttempt", "target_team_id", "team(id)"},
	}
}
