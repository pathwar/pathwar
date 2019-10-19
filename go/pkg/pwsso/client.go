package pwsso

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type Opts struct {
	AllowUnsafe bool
	Logger      *zap.Logger
	ClientID    string
	// ClientSecret string
}

type Client interface {
	TokenWithClaims(bearer string) (*jwt.Token, jwt.MapClaims, error)
	Whoami(token string) (map[string]interface{}, error)
	Logout(token string) error
}

type client struct {
	publicKey interface{} // result of x509.ParsePKIXPublicKey
	logger    *zap.Logger
	realm     string
	clientID  string
	opts      Opts
}

func New(publicKey string, realm string, opts Opts) (Client, error) {
	c := &client{
		opts:   opts,
		realm:  realm,
		logger: opts.Logger,
	}

	if c.opts.ClientID == "" {
		c.opts.ClientID = "platform-cli"
	}

	{ // parse public key
		key := []byte(fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----\n", publicKey))
		pubPem, _ := pem.Decode(key)
		if pubPem == nil {
			return nil, errors.New("invalid pubkey")
		}
		parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse public key: %w", err)
		}
		c.publicKey = parsedKey
	}

	return c, nil
}
