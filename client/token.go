package client

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	fmt "fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

func TokenWithClaims(bearer string, opts Options) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	// FIXME: add an option to automatically fetch the public key from
	// https://id.pathwar.land/auth/realms/Pathwar-Dev/protocol/openid-connect/certs
	// or
	// https://id.pathwar.land/auth/realms/Pathwar-Dev
	token, err := jwt.ParseWithClaims(bearer, claims, keyFunc(opts))
	if err != nil {
		if opts.Unsafe {
			zap.L().Warn(
				"invalid token",
				zap.Error(err),
				zap.Bool("client-unsafe", true),
			)
			parser := new(jwt.Parser)
			token, _, err := parser.ParseUnverified(bearer, claims)
			return token, claims, err
		}
		return nil, nil, err
	}
	return token, claims, nil
}

func keyFunc(opts Options) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		key := []byte(fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----\n", opts.PublicKey))
		pubPem, _ := pem.Decode(key)
		if pubPem == nil {
			return nil, errors.New("invalid pubkey")
		}
		parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
		return parsedKey, err
	}
}
