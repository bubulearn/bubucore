package bubucore

import (
	"github.com/bubulearn/bubucore/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

// ParseAccessToken parses access token and returns its claims
func ParseAccessToken(tokenContent string) (*AccessTokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenContent, &AccessTokenClaims{}, parseJWTKeyFunc)
	if err != nil {
		logrus.Warn("failed to parse JWT (access token): ", tokenContent, ": ", err)
		return nil, ErrTokenInvalid
	}
	claims := parsed.Claims.(*AccessTokenClaims)
	err = claims.Valid()
	if err != nil {
		logrus.Warn("failed to validate JWT (access token): ", tokenContent, ": ", err)
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// parseJWTKeyFunc returns JWT password key
func parseJWTKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrTokenUnsupported
	}
	return Opt.JWTPassword, nil
}

// TokenClaims is token claims interface
type TokenClaims interface {
	GetTokenID() string
	GetUserID() string
	GetRelatedTokenID() string
	jwt.Claims
}

// TokenClaimsDft is a default realization of TokenClaims
type TokenClaimsDft struct {
	UserID string `json:"uid"`
	jwt.StandardClaims
}

// GetTokenID returns token ID
func (c TokenClaimsDft) GetTokenID() string {
	return c.Id
}

// GetUserID returns token's user ID
func (c TokenClaimsDft) GetUserID() string {
	return c.UserID
}

// Valid checks is data in claims is valid
func (c TokenClaimsDft) Valid() error {
	if !utils.ValidateUUID(c.GetUserID()) {
		return ErrTokenInvalid
	}
	if !utils.ValidateUUID(c.GetTokenID()) {
		return ErrTokenInvalid
	}

	now := time.Now().Unix()

	if !c.VerifyExpiresAt(now, true) {
		return ErrTokenExpired
	}
	if !c.VerifyIssuedAt(now, false) {
		return ErrTokenInvalid
	}
	if c.VerifyNotBefore(now, false) == false {
		return ErrTokenInvalid
	}

	return nil
}

// GetRelatedTokenID returns related refresh token ID
func (c TokenClaimsDft) GetRelatedTokenID() string {
	return ""
}

// AccessTokenClaims is an access token claims
type AccessTokenClaims struct {
	TokenClaimsDft
	Role           int    `json:"rl"`
	Name           string `json:"nm,omitempty"`
	RefreshTokenID string `json:"rti,omitempty"`
}

// Valid checks is data in claims is valid
func (c AccessTokenClaims) Valid() error {
	if c.Role == 0 {
		return ErrTokenInvalid
	}
	if c.GetRelatedTokenID() != "" && !utils.ValidateUUID(c.GetRelatedTokenID()) {
		return ErrTokenInvalid
	}
	return c.TokenClaimsDft.Valid()
}

// GetRelatedTokenID returns related refresh token ID
func (c AccessTokenClaims) GetRelatedTokenID() string {
	return c.RefreshTokenID
}
