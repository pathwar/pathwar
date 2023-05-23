package pwsso

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

type Opts struct {
	AllowUnsafe bool
	Logger      *zap.Logger
	ClientID    string

	// following fields are not (yet) used by this package but are used to configure other SSO related stuff

	ClientSecret string
	Realm        string
	TokenFile    string
	Pubkey       string
}

// NewOpts returns sane default values for development
func NewOpts() Opts {
	return Opts{
		Pubkey:       "",
		ClientID:     testingClientID,
		Realm:        testingRealm,
		ClientSecret: "",
		TokenFile:    "default",
		AllowUnsafe:  false,
		Logger:       zap.NewNop(),
	}
}

func (opts *Opts) ApplyDefaults() {
	if opts.Pubkey == "" {
		opts.Pubkey = testingPubKey
	}
}

type Client interface {
	TokenWithClaims(bearer string) (*jwt.Token, jwt.MapClaims, error)
	Whoami(token string) (map[string]interface{}, error)
}

type client struct {
	publicKey interface{} // result of x509.ParsePKIXPublicKey
	logger    *zap.Logger
	realm     string
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
			return nil, errcode.ErrSSOInvalidPublicKey
		}

		parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
		if err != nil {
			return nil, errcode.ErrSSOInvalidPublicKey.Wrap(err)
		}
		c.publicKey = parsedKey
	}

	return c, nil
}
