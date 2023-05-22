package pwsso

import (
	"net/http"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"moul.io/roundtripper"
)

const (
	testingPubKey   = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzvJIRVk38uby5bGVJpCS\ndr7RardC6s1G61A+CO127rvTjLkDFxhM3n6NkF1GVBvIXBbvaj6Q7+CKPR1L5NMG\nlEFvQTjbcuBL11v7ViYE+UnNmJcXHZb+kzVml3evcUhyVqf8aPkyT7CgzM+0BjPf\nUYFZ4raWM9vG+WAmYCXMnKek4jFhLhZGQO9n9W7wrZW3Yegc/YQWuqGtkaRUsfwd\nwQJn4OIhpMVw4YKQIpz7BPObRqAh49dn1waQ5TEvW0IUVwHW8nTCHbePXxLeSEat\n0REs32wJt5G9JgSnaqs/j7AqctG41qbO0dqxE/FgmcAsCmd82MUFI1VBzOYmnLdT\njQIDAQAB"
	testingRealm    = "undefined"
	testingClientID = "bJpLWOLTRseEVfM9kvFhKfi9wUBmm8Gh"
	testingToken    = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlBPNVlab3NHZ0V1VFViQklfNEFpMiJ9.eyJlbWFpbCI6Im1pa2FlbEBiZXJ0eS50ZWNoIiwibmlja25hbWUiOiJtaWthZWwiLCJjb2xvciI6ImJsdWUiLCJpc3MiOiJodHRwczovL2Rldi01Y2N3enk4cXRjc2pzbnBmLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJnb29nbGUtb2F1dGgyfDExNTkwMDE5MDAxMDA2MjA1OTI0MiIsImF1ZCI6WyJodHRwczovL3BhdGh3YXIubmV0LyIsImh0dHBzOi8vZGV2LTVjY3d6eThxdGNzanNucGYudXMuYXV0aDAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTY4NDUxMDkwMiwiZXhwIjoxNjg0NTk3MzAyLCJhenAiOiJiSnBMV09MVFJzZUVWZk05a3ZGaEtmaTl3VUJtbThHaCIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJwZXJtaXNzaW9ucyI6WyJhZG1pbiIsImFnZW50Il19.X5zm6w2CNX2DnoBOcgPS8zsFuT1yJb3tHmCY65sPucGPYtyMLQV7fsVLUp-Rxq68iBn1kwRSZ9vOH5cYM7yqpou44w3fgaeFY1_oWBPw3xtGD_mCghIpORQ9BAQv6GHgkM7di0Q9I2gCMP8sMAYOBrxp43GeKNsDFuQzD2bj-rxwEaXhN3gPV8hznsuXpT0BfwpQNSaEJY_pWYkwvPncjrDStfMmdqkG5EV5AD02kY3h1umKJ1DbA1Vsj8WSnA2_ae7eGQ9onSLBKv7th_D_5RtVrM-H2tkrcCmLOAbhuQVfdOEqesrR1dUMOQ5D3sRsMSSpkFx7jHCHNsH5kSITUg"
	testingToken2   = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlBPNVlab3NHZ0V1VFViQklfNEFpMiJ9.eyJlbWFpbCI6Im1pa2FlbC52YWxsZW5ldHByb0BnbWFpbC5jb20iLCJuaWNrbmFtZSI6Ik1pa2F0ZWNoIiwiY29sb3IiOiJibHVlIiwiaXNzIjoiaHR0cHM6Ly9kZXYtNWNjd3p5OHF0Y3Nqc25wZi51cy5hdXRoMC5jb20vIiwic3ViIjoiZ2l0aHVifDcxMzM5MTUzIiwiYXVkIjpbImh0dHBzOi8vcGF0aHdhci5uZXQvIiwiaHR0cHM6Ly9kZXYtNWNjd3p5OHF0Y3Nqc25wZi51cy5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNjg0NTA5NjIzLCJleHAiOjE2ODQ1OTYwMjMsImF6cCI6ImJKcExXT0xUUnNlRVZmTTlrdkZoS2ZpOXdVQm1tOEdoIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBvZmZsaW5lX2FjY2VzcyIsInBlcm1pc3Npb25zIjpbXX0.vaoySqu8G5wfdSRiB9T6v2HKS3deGk7QyCGWB08G0kUjhwXYS7igk0ajSn0YhlVQqOMb5zmk7zbmu5WWBdi7tEihAZ-RXQ2fn-L9lhPkYiDexpOIJfb0XdMapmDx8oCRzqCri2rrE6S_eUGpZm5mwOZL2WarABRj7q_-vttiQhHFo7isbUOOx7FGU-GlzxZPe42inOmo0Cl5RL9S6lvfWRPGZFnSr-Mco5RAmaL6d6r1Z7gAy0KxHuHmo4dSbKM2KemDOUgqnj2ngReRXrX3RAI5AdLQlqLgii-3QFeGKuISwHpvpQ8vVpgoqmCViH9o0F3ioVpQSnHJG7cQ8M6PCQ"
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

func TestingToken2(t *testing.T) *jwt.Token {
	t.Helper()
	token, _, err := TokenWithClaims(testingToken2, testingPubKey, true)
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
		Realm:       testingRealm,
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
