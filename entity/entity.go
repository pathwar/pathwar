package entity

func ByName(name string) interface{} {
	return AllMap()[name]
}

func AllMap() map[string]interface{} {
	return map[string]interface{}{
		"Achievement":       Achievement{},
		"AuthMethod":        AuthMethod{},
		"Coupon":            Coupon{},
		"CouponValidation":  CouponValidation{},
		"Event":             Event{},
		"Hypervisor":        Hypervisor{},
		"InventoryItem":     InventoryItem{},
		"Level":             Level{},
		"LevelFlavor":       LevelFlavor{},
		"LevelInstance":     LevelInstance{},
		"LevelSubscription": LevelSubscription{},
		"LevelValidation":   LevelValidation{},
		"LevelVersion":      LevelVersion{},
		"Notification":      Notification{},
		"Team":              Team{},
		"TeamMember":        TeamMember{},
		"Tournament":        Tournament{},
		"TournamentMember":  TournamentMember{},
		"TournamentTeam":    TournamentTeam{},
		"User":              User{},
		"UserSession":       UserSession{},
		"WhoswhoAttempt":    WhoswhoAttempt{},
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
		{"TeamMember", "user_id", "user(id)"},
		{"TeamMember", "team_id", "team(id)"},
		{"TournamentMember", "user_id", "user(id)"},
		{"TournamentMember", "tournament_team_id", "tournament_team(id)"},
		{"UserSession", "user_id", "user(id)"},
		{"LevelVersion", "level_id", "level(id)"},
		{"TournamentTeam", "tournament_id", "tournament(id)"},
		{"TournamentTeam", "team_id", "team(id)"},
		{"WhoswhoAttempt", "author_id", "tournament_member(id)"},
		{"WhoswhoAttempt", "target_member_id", "tournament_member(id)"},
		{"WhoswhoAttempt", "target_team_id", "tournament_team(id)"},
		{"LevelValidation", "level_subscription_id", "level_subscription(id)"},
		{"LevelSubscription", "tournament_team_id", "tournament_team(id)"},
		{"LevelSubscription", "level_flavor_id", "level_flavor(id)"},
		{"InventoryItem", "owner_id", "tournament_member(id)"},
		{"LevelFlavor", "level_id", "level(id)"},
		{"LevelInstance", "hypervisor_id", "hypervisor(id)"},
		{"LevelInstance", "flavor_id", "level_flavor(id)"},
		{"Notification", "user_id", "user(id)"},
		{"Coupon", "tournament_id", "tournament(id)"},
		{"CouponValidation", "author_id", "tournament_member(id)"},
		{"CouponValidation", "coupon_id", "coupon(id)"},
		{"Achievement", "author_id", "tournament_member(id)"},
		{"Achievement", "level_validation_id", "level_validation(id)"},
	}
}
