package entity

func All() []interface{} {
	return []interface{}{
		Level{},
		UserSession{},
		User{},
		TeamMember{},
		Team{},
	}
}
