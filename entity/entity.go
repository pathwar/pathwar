package entity

func ByName(name string) interface{} {
	return AllMap()[name]
}

func AllMap() map[string]interface{} {
	return map[string]interface{}{
		"Achievement":       Achievement{},
		"Coupon":            Coupon{},
		"Event":             Event{},
		"Hypervisor":        Hypervisor{},
		"LevelFlavor":       LevelFlavor{},
		"LevelInstance":     LevelInstance{},
		"LevelSubscription": LevelSubscription{},
		"Level":             Level{},
		"Notification":      Notification{},
		"ShopItem":          ShopItem{},
		"TeamMember":        TeamMember{},
		"Team":              Team{},
		"TournamentTeam":    TournamentTeam{},
		"Tournament":        Tournament{},
		"UserSession":       UserSession{},
		"User":              User{},
		"WhoswhoAttempt":    WhoswhoAttempt{},
	}
}

func All() []interface{} {
	return []interface{}{
		Achievement{},
		Coupon{},
		Event{},
		Hypervisor{},
		LevelFlavor{},
		LevelInstance{},
		LevelSubscription{},
		Level{},
		Notification{},
		ShopItem{},
		TeamMember{},
		Team{},
		TournamentTeam{},
		Tournament{},
		UserSession{},
		User{},
		WhoswhoAttempt{},
	}
}
