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
		"ChallengeVersion":      ChallengeVersion{},
		"Coupon":                Coupon{},
		"CouponValidation":      CouponValidation{},
		"Hypervisor":            Hypervisor{},
		"InventoryItem":         InventoryItem{},
		"Notification":          Notification{},
		"Organization":          Organization{},
		"OrganizationMember":    OrganizationMember{},
		"Season":                Season{},
		"TeamMember":            TeamMember{},
		"Team":                  Team{},
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
		{"Achievement", "author_id", "team_member(id)"},
		{"Achievement", "challenge_validation_id", "challenge_validation(id)"},
		{"ChallengeFlavor", "challenge_version_id", "challenge_version(id)"},
		{"ChallengeInstance", "flavor_id", "challenge_flavor(id)"},
		{"ChallengeInstance", "hypervisor_id", "hypervisor(id)"},
		{"ChallengeSubscription", "challenge_flavor_id", "challenge_flavor(id)"},
		{"ChallengeSubscription", "team_id", "team(id)"},
		{"ChallengeValidation", "challenge_subscription_id", "challenge_subscription(id)"},
		{"ChallengeValidation", "team_member_id", "team_member(id)"},
		{"ChallengeVersion", "challenge_id", "challenge(id)"},
		{"Coupon", "season_id", "season(id)"},
		{"CouponValidation", "author_id", "team_member(id)"},
		{"CouponValidation", "coupon_id", "coupon(id)"},
		{"InventoryItem", "owner_id", "team_member(id)"},
		{"Notification", "user_id", "user(id)"},
		{"OrganizationMember", "organization_id", "organization(id)"},
		{"OrganizationMember", "user_id", "user(id)"},
		{"TeamMember", "team_id", "team(id)"},
		{"TeamMember", "user_id", "user(id)"},
		{"Team", "organization_id", "organization(id)"},
		{"Team", "season_id", "season(id)"},
		{"WhoswhoAttempt", "author_id", "team_member(id)"},
		{"WhoswhoAttempt", "target_member_id", "team_member(id)"},
		{"WhoswhoAttempt", "target_organization_id", "team(id)"},
	}
}
