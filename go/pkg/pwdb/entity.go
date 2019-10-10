package pwdb

func ByName(name string) interface{} {
	return AllMap()[name]
}

func AllMap() map[string]interface{} {
	return map[string]interface{}{
		"Achievement":           Achievement{},
		"Coupon":                Coupon{},
		"CouponValidation":      CouponValidation{},
		"Hypervisor":            Hypervisor{},
		"InventoryItem":         InventoryItem{},
		"Challenge":             Challenge{},
		"ChallengeFlavor":       ChallengeFlavor{},
		"ChallengeInstance":     ChallengeInstance{},
		"ChallengeSubscription": ChallengeSubscription{},
		"ChallengeValidation":   ChallengeValidation{},
		"ChallengeVersion":      ChallengeVersion{},
		"Notification":          Notification{},
		"Team":                  Team{},
		"TeamMember":            TeamMember{},
		"Tournament":            Tournament{},
		"TournamentMember":      TournamentMember{},
		"TournamentTeam":        TournamentTeam{},
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
		{"Achievement", "author_id", "tournament_member(id)"},
		{"Achievement", "challenge_validation_id", "challenge_validation(id)"},
		{"Coupon", "tournament_id", "tournament(id)"},
		{"CouponValidation", "author_id", "tournament_member(id)"},
		{"CouponValidation", "coupon_id", "coupon(id)"},
		{"InventoryItem", "owner_id", "tournament_member(id)"},
		{"ChallengeFlavor", "challenge_version_id", "challenge_version(id)"},
		{"ChallengeInstance", "flavor_id", "challenge_flavor(id)"},
		{"ChallengeInstance", "hypervisor_id", "hypervisor(id)"},
		{"ChallengeSubscription", "challenge_flavor_id", "challenge_flavor(id)"},
		{"ChallengeSubscription", "tournament_team_id", "tournament_team(id)"},
		{"ChallengeValidation", "challenge_subscription_id", "challenge_subscription(id)"},
		{"ChallengeValidation", "tournament_member_id", "tournament_member(id)"},
		{"ChallengeVersion", "challenge_id", "challenge(id)"},
		{"Notification", "user_id", "user(id)"},
		{"TeamMember", "team_id", "team(id)"},
		{"TeamMember", "user_id", "user(id)"},
		{"TournamentMember", "tournament_team_id", "tournament_team(id)"},
		{"TournamentMember", "user_id", "user(id)"},
		{"TournamentTeam", "team_id", "team(id)"},
		{"TournamentTeam", "tournament_id", "tournament(id)"},
		// {"User", "active_tournament_member_id", "tournament_member(id)"}, // FIXME: check why this cause an error!?
		{"WhoswhoAttempt", "author_id", "tournament_member(id)"},
		{"WhoswhoAttempt", "target_member_id", "tournament_member(id)"},
		{"WhoswhoAttempt", "target_team_id", "tournament_team(id)"},
	}
}
