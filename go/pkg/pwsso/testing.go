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
	testingRealm    = "Pathwar-Dev"
	testingClientID = "bJpLWOLTRseEVfM9kvFhKfi9wUBmm8Gh"
	testingToken    = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlBPNVlab3NHZ0V1VFViQklfNEFpMiJ9.eyJpc3MiOiJodHRwczovL2Rldi01Y2N3enk4cXRjc2pzbnBmLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJnaXRodWJ8NzEzMzkxNTMiLCJhdWQiOiJodHRwczovL3BhdGh3YXIubmV0LyIsImlhdCI6MTY4NDM3NzM4MiwiZXhwIjoxNjg0NDYzNzgyLCJhenAiOiJiSnBMV09MVFJzZUVWZk05a3ZGaEtmaTl3VUJtbThHaCIsInNjb3BlIjoib2ZmbGluZV9hY2Nlc3MiLCJwZXJtaXNzaW9ucyI6WyJhZG1pbiIsImFnZW50Il19.VctZTCvkSPGdHN5h3WjveCOPIxHtGIK2yam0PGHJSiWzTUTGxy4PM_0t5u2IywI2gLv6YSdxH9qj9FEuEWO37IeBHRc3n1lT_15MqMha__Zk5Ps-C4uHnCiyJQD23m1Zb-eupmjTubCJ5ua1nYmQ_eY9-YUhhnU9CsZeW5S0feEbmIS7bHLmduPV-iqLRuCiqEdk8y0QAQjwZ050SKOTzyIkJFHzG3b8909cFc48EAfyRdzEeXEY0x8au2B87dicPQtDh1Sb_c0_UVh2s8xVadiciIx11hr8bTgltlkfeNvWrByYNMgTlFR8btZc6sGu4M0xGqjLkLV129HQmiN3dA"
	testingToken2   = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlBPNVlab3NHZ0V1VFViQklfNEFpMiJ9.eyJpc3MiOiJodHRwczovL2Rldi01Y2N3enk4cXRjc2pzbnBmLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJnb29nbGUtb2F1dGgyfDExNTkwMDE5MDAxMDA2MjA1OTI0MiIsImF1ZCI6Imh0dHBzOi8vcGF0aHdhci5uZXQvIiwiaWF0IjoxNjg0NDg0NDAxLCJleHAiOjE2ODQ1NzA4MDEsImF6cCI6ImJKcExXT0xUUnNlRVZmTTlrdkZoS2ZpOXdVQm1tOEdoIiwic2NvcGUiOiJvZmZsaW5lX2FjY2VzcyIsInBlcm1pc3Npb25zIjpbXX0.Fv5iVgZYc3OHSoBs38hYhDBD471TPJ2nyOtyp-JsNGSsXd4uoh84XYqRWF2WUFyNGw3pnJZyyr1gO4GYo7BXunfaLTYBzHK8CZ5dJnsCVqH342WdWYOaYpouBCyvofq3k3cdu5mjvTDyswV9logJlAtHuNAqMKSoxU0hU4_BdG6X5-TmMmVQkw0CGpcfcv7PYfeWWB0nyCETkPew-LZsSchRh2h2aLCe3XOUK0vhsuyAhku3HgG1rLmEubxGmGTieBDpf0KyWPIH6Kg4EA_Ni6pQ40nbpkQvPbLVNeBXKdkvym658zeygADV048LXQMRQnez7RqwxLc0VePxtAOc7Q"
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
