package pwsso

import (
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

const (
	testingPubKey   = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlEFxLlywsbI5BQ7DVkA66fICWGIYPpD+aZNYRR7SIc0zdtJR4xMOt5CjM0vbYT4z2a1U2yl0ewunyxFm8niS8w6mKYFnOS4nnSchQyIAmJkpLC4eAjijCdEHdr8mSqamThSrVRGSYEEsa+adidC13kRDy7NDKhvZb8F0YqnktNk6WHSlb8r2QRLPJ1DX534jjXPY6l/eoHuLJAOZxBlfwV5Dg37TVmf2xAH812E7ZigycLAvhsMvr5x2jLavAEEnZZmlQf4cyQ4tlMzKS1Zp0NcdOGS/i6lrndc5pNtZQuGr8IGBrEbTRFUiavn/HDnyalYZy8T5LakXRdVaKdshAQIDAQAB"
	testingRealm    = "Pathwar-Dev"
	testingClientID = "platform-cli"
	testingToken    = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJDck10ZmN1cjFDcVNtT28teHZacUt0ZTRoODk4ZjZpYl9KOGk5TXZDck5zIn0.eyJqdGkiOiI2Y2RmYzNiNy1lYjZhLTQ2YjktOWU1OC03YzhjMTVmOTRlNzAiLCJleHAiOjE1Njc2ODU0NTgsIm5iZiI6MCwiaWF0IjoxNTY3Njg1MTU4LCJpc3MiOiJodHRwczovL2lkLnBhdGh3YXIubGFuZC9hdXRoL3JlYWxtcy9QYXRod2FyLURldiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIwNDgyNjZiOS0yY2M4LTQ2ZjMtOTcyZC0zN2YyZDhmY2M3NWIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJwbGF0Zm9ybS1mcm9udCIsIm5vbmNlIjoiOWZjOGNjYTktNDMyMy00YjA3LTg3NmUtNzE2NTcyN2M5NDc3IiwiYXV0aF90aW1lIjoxNTY3Njc3NTMzLCJzZXNzaW9uX3N0YXRlIjoiN2I3OWI1NmUtZjUzNS00MmE0LTliYWUtYzk2YjE1NGMxNWZjIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWFuZnJlZCBUb3Vyb24iLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJtb3VsIiwiZ2l2ZW5fbmFtZSI6Ik1hbmZyZWQiLCJmYW1pbHlfbmFtZSI6IlRvdXJvbiIsImVtYWlsIjoibUA0Mi5hbSJ9.YdjFGsdvnWFTC6ORhb8jXMzFjGMEttdNLIHuc8OMRYW5UUi4GleOu81IyPmj9GQdIsoew7KHFd5y0LpaoaXpwsSVUDhOHiBQWD11xyk8X2ULCYdcHZB4jyIeb9TftdYDyZmtdhNvZXXSod6pL71JKUA2BcMYrxNrp98qnT2-pUCPFGZ85Lcdz0MiAgldaX_rDiS7JlSZW5McJdNq--JZd9hZA1pomds_yMr3N9kvXnMmgIgVsAKmrvXmOJ9sXDJlFWV7yGb29ZFf1sE1meEDPlOj2ZMu5X6NInXbCl_dWKPgmPf-qELfu33LOeDTVdp-E7NaHyj2rduN4gOSH3H_Ag"
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
