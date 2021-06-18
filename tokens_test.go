package bubucore_test

import (
	"github.com/bubulearn/bubucore"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenClaims_Valid(t *testing.T) {
	valid := &bubucore.AccessTokenClaims{
		RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		Role:           10,
		Name:           "User Name",
	}
	valid.TokenClaimsDft = bubucore.TokenClaimsDft{
		UserID: "03a4e59c-fb22-4bfa-8739-8062bcdd2005",
		StandardClaims: jwt.StandardClaims{
			Id:        "0bf97df4-6246-4809-bdf7-e8d993668283",
			ExpiresAt: time.Now().Unix() + int64(10*time.Minute.Seconds()),
		},
	}

	invalid := []*bubucore.AccessTokenClaims{
		{},
		{Role: 10, RefreshTokenID: ""},
		{
			bubucore.TokenClaimsDft{},
			20,
			"User",
			"145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
		{
			TokenClaimsDft: bubucore.TokenClaimsDft{
				UserID: "5881a508-65b3-4caf-930b-07e7c363d7e2",
			},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
		{
			TokenClaimsDft: bubucore.TokenClaimsDft{
				UserID: "5881a508-65b3-4caf-930b-07e7c363d7e2",
				StandardClaims: jwt.StandardClaims{
					Id:        "dd5d2ced-e168-4794-822e-f13fe952dddd",
					ExpiresAt: 100,
				},
			},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
		{
			TokenClaimsDft: bubucore.TokenClaimsDft{
				UserID: "5881a508-65b3-4caf-930b-07e7c363d7e2",
				StandardClaims: jwt.StandardClaims{
					Id:        "dd5d2ced-e168-4794-822e-f13fe952dddd",
					ExpiresAt: time.Now().Unix() + int64(10*time.Minute.Seconds()),
					NotBefore: time.Now().Unix() + int64(5*time.Minute.Seconds()),
				},
			},
			Role:           20,
			Name:           "User",
			RefreshTokenID: "145b9a16-9a57-4264-b7a7-b96ab3b7e7b9",
		},
	}

	assert.NoError(t, valid.Valid())

	for _, claims := range invalid {
		assert.Error(t, claims.Valid())
	}
}
