package pwsso

import (
	"net/http"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"moul.io/roundtripper"
)

const (
	testingPubKey   = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlEFxLlywsbI5BQ7DVkA66fICWGIYPpD+aZNYRR7SIc0zdtJR4xMOt5CjM0vbYT4z2a1U2yl0ewunyxFm8niS8w6mKYFnOS4nnSchQyIAmJkpLC4eAjijCdEHdr8mSqamThSrVRGSYEEsa+adidC13kRDy7NDKhvZb8F0YqnktNk6WHSlb8r2QRLPJ1DX534jjXPY6l/eoHuLJAOZxBlfwV5Dg37TVmf2xAH812E7ZigycLAvhsMvr5x2jLavAEEnZZmlQf4cyQ4tlMzKS1Zp0NcdOGS/i6lrndc5pNtZQuGr8IGBrEbTRFUiavn/HDnyalYZy8T5LakXRdVaKdshAQIDAQAB"
	testingRealm    = "Pathwar-Dev"
	testingClientID = "platform-cli"
	testingToken    = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJDck10ZmN1cjFDcVNtT28teHZacUt0ZTRoODk4ZjZpYl9KOGk5TXZDck5zIn0.eyJqdGkiOiI0ZGE4ZTM2NS1iZTkzLTRmMGEtYmU0ZC0yNDdjMzA4OGZmNWUiLCJleHAiOjE1ODM0Mjc1MTIsIm5iZiI6MCwiaWF0IjoxNTgzNDI3MjEyLCJpc3MiOiJodHRwczovL2lkLnBhdGh3YXIubGFuZC9hdXRoL3JlYWxtcy9QYXRod2FyLURldiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIwNDgyNjZiOS0yY2M4LTQ2ZjMtOTcyZC0zN2YyZDhmY2M3NWIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJwbGF0Zm9ybS1jbGkiLCJub25jZSI6IjZlMWUyYjc4LTk0MjgtNDRhNi04ZjIwLTA5NTY3ZTE1Y2FjMyIsImF1dGhfdGltZSI6MTU4MzQyNzIwNywic2Vzc2lvbl9zdGF0ZSI6ImEyMzg2N2U2LTc0ZjEtNGRmOS04ZDRiLWU5NTVlMWRmMmYxNCIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsicGxhdGZvcm0tY2xpIjp7InJvbGVzIjpbImFnZW50IiwiYWRtaW4iXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSBvZmZsaW5lX2FjY2VzcyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWFuZnJlZCBUb3Vyb24iLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJtb3VsIiwiZ2l2ZW5fbmFtZSI6Ik1hbmZyZWQiLCJmYW1pbHlfbmFtZSI6IlRvdXJvbiIsImVtYWlsIjoibUA0Mi5hbSJ9.I9jYiBGCacaBiqndq1EsinZxY-uWRjdHZbFRdE9CWsSiOEJzKGznufEppk0bj2XmAm4GwfWey55U-jHh91KgnDJG7XsgA2p_t-LX1yj4EgrHxcXQ0PiOKU19br4kbCfKVaOMsBQqa-pGyZVFwVc9rYmGA6xtx6No1O5j-tdsizp5-HVNil0E195ZnSoMiNk9yJsG8-ta7wrQ6u9PqPbnEuhltu6SZyfAD7gTw2RUDu77LKISIaJCPbD5IPj2Rtv2gfM4BoZ8TiMYO_DSRIAWsFc1C1z8iR6-BvAvOAfqDV4GeyD9DQsMDxz5qYmTnHnXMrVNSvYd6aehwyDik-ERIA"
)

func TestingClaims(t *testing.T) *Claims {
	t.Helper()

	token := TestingToken(t)
	return ClaimsFromToken(token)
}

func TestingToken(t *testing.T) *jwt.Token {
	t.Helper()
	token, _, err := TokenWithClaims(testingToken, testingPubKey, true)
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
