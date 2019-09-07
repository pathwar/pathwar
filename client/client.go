package client

import (
	jwt "github.com/dgrijalva/jwt-go"
	"pathwar.land/entity"
)

func UserSessionFromToken(token *jwt.Token) (entity.UserSession, error) {
	claims := token.Claims.(jwt.MapClaims)
	sess := entity.UserSession{
		// FIXME: use mapstructure
		Username: claims["preferred_username"].(string),
	}
	return sess, nil
}
