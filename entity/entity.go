package entity

func All() []interface{} {
	return []interface{}{
		Level{},
		UserSession{},
	}
}
