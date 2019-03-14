package entity

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
