package pwsso

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"moul.io/roundtripper"
)

const (
	testingPubKey   = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzvJIRVk38uby5bGVJpCS\ndr7RardC6s1G61A+CO127rvTjLkDFxhM3n6NkF1GVBvIXBbvaj6Q7+CKPR1L5NMG\nlEFvQTjbcuBL11v7ViYE+UnNmJcXHZb+kzVml3evcUhyVqf8aPkyT7CgzM+0BjPf\nUYFZ4raWM9vG+WAmYCXMnKek4jFhLhZGQO9n9W7wrZW3Yegc/YQWuqGtkaRUsfwd\nwQJn4OIhpMVw4YKQIpz7BPObRqAh49dn1waQ5TEvW0IUVwHW8nTCHbePXxLeSEat\n0REs32wJt5G9JgSnaqs/j7AqctG41qbO0dqxE/FgmcAsCmd82MUFI1VBzOYmnLdT\njQIDAQAB"
	testingRealm    = "Pathwar-Dev"
	testingClientID = "bJpLWOLTRseEVfM9kvFhKfi9wUBmm8Gh"
	testingToken    = "x_TrDAWz47HbRJt5ltiOqB8y15gYjCfBgsBg_RzLztFuW"
	testingToken2   = "EqHeaU05F_pK7qNoMFkRGnC_uRGj0pg7L-n61s6iieK9k"
)

func TestingClaims(t *testing.T) *Claims {
	t.Helper()

	token := TestingToken(t)
	return ClaimsFromToken(token)
}

func TestingToken(t *testing.T) *jwt.Token {
	t.Helper()
	accesToken, err := GetTestingTokenFromRefresh(testingToken)
	if err != nil {
		t.Fatalf("get token from refresh: %v", err)
	}
	token, _, err := TokenWithClaims(accesToken, testingPubKey, true)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	return token
}

func TestingToken2(t *testing.T) *jwt.Token {
	t.Helper()
	accesToken, err := GetTestingTokenFromRefresh(testingToken2)
	if err != nil {
		t.Fatalf("get token from refresh: %v", err)
	}
	token, _, err := TokenWithClaims(accesToken, testingPubKey, true)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	return token
}

func TestingSSO(t *testing.T, logger *zap.Logger) Client {
	t.Helper()
	ssoOpts := Opts{
		AllowUnsafe: true,
		Logger:      logger,
		ClientID:    testingClientID,
	}
	sso, err := New(testingPubKey, testingRealm, ssoOpts)
	if err != nil {
		t.Fatalf("init SSO: %v", err)
	}

	return sso
}

func TestingTransport(t *testing.T) http.RoundTripper {
	return &roundtripper.Transport{
		ExtraHeader: http.Header{
			"Authorization": []string{"Bearer " + testingToken},
		},
	}
}

func GetTestingTokenFromRefresh(token string) (string, error) {
	values := map[string]string{"grant_type": "refresh_token", "client_id": testingClientID, "refresh_token": token}
	jsonData, err := json.Marshal(values)
	resp, err := http.Post(ProviderTokenURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	return res["access_token"].(string), nil
}
