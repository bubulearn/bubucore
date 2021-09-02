package tokens

import (
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/utils"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"time"
)

// ParseAccessToken parses access token and returns its claims
func ParseAccessToken(tokenContent string) (*AccessTokenClaims, error) {
	parser := &jwt.Parser{
		ValidMethods: []string{
			jwt.SigningMethodHS256.Alg(),
			jwt.SigningMethodHS384.Alg(),
			jwt.SigningMethodHS512.Alg(),
		},
	}
	parsed, err := parser.ParseWithClaims(tokenContent, &AccessTokenClaims{}, parseJWTKeyFunc)
	if err != nil {
		log.Warn("failed to parse JWT (access token): ", tokenContent, ": ", err)
		return nil, bubucore.ErrTokenInvalid
	}
	claims := parsed.Claims.(*AccessTokenClaims)
	return claims, nil
}

// parseJWTKeyFunc returns JWT password key
func parseJWTKeyFunc(_ *jwt.Token) (interface{}, error) {
	return bubucore.Opt.JWTPassword, nil
}

// TokenClaims is token claims interface
type TokenClaims interface {
	// GetTokenID returns token ID
	GetTokenID() string

	// GetUserID returns token's user ID
	GetUserID() string

	// GetRelatedTokenID returns related refresh token ID
	GetRelatedTokenID() string

	// GetAllowedServices returns services allowed to use token with.
	// Return nil or an empty slice if all services are allowed.
	GetAllowedServices() []string

	jwt.Claims
}

// TokenClaimsDft is a default realization of TokenClaims
type TokenClaimsDft struct {
	UserID          string   `json:"uid"`
	ServicesAllowed []string `json:"srvs,omitempty"`
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

// GetAllowedServices returns services allowed to use token with.
// Return nil or an empty slice if all services are allowed.
func (c TokenClaimsDft) GetAllowedServices() []string {
	return c.ServicesAllowed
}

// Valid checks is data in claims is valid
func (c TokenClaimsDft) Valid() error {
	if !utils.ValidateUUID(c.GetUserID()) {
		return bubucore.ErrTokenInvalid
	}
	if !utils.ValidateUUID(c.GetTokenID()) {
		return bubucore.ErrTokenInvalid
	}

	now := time.Now().Unix()

	if !c.VerifyExpiresAt(now, true) {
		return bubucore.ErrTokenExpired
	}
	if !c.VerifyIssuedAt(now, false) {
		return bubucore.ErrTokenInvalid
	}
	if c.VerifyNotBefore(now, false) == false {
		return bubucore.ErrTokenInvalid
	}

	allowed := c.GetAllowedServices()
	if len(allowed) > 0 {
		currSrv := bubucore.Opt.ServiceName
		if currSrv == "" {
			return bubucore.ErrTokenInvalid
		}
		ok := false
		for _, s := range allowed {
			if s == currSrv {
				ok = true
				break
			}
		}
		if !ok {
			return bubucore.ErrTokenInvalid
		}
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
	Role           int           `json:"rl"`
	Name           string        `json:"nm,omitempty"`
	RefreshTokenID string        `json:"rti,omitempty"`
	Language       i18n.Language `json:"lng,omitempty"`
}

// Valid checks is data in claims is valid
func (c AccessTokenClaims) Valid() error {
	if c.Role == 0 {
		return bubucore.ErrTokenInvalid
	}
	if c.GetRelatedTokenID() != "" && !utils.ValidateUUID(c.GetRelatedTokenID()) {
		return bubucore.ErrTokenInvalid
	}
	return c.TokenClaimsDft.Valid()
}

// GetRelatedTokenID returns related refresh token ID
func (c AccessTokenClaims) GetRelatedTokenID() string {
	return c.RefreshTokenID
}

// RefreshTokenClaims is an refresh token claims
type RefreshTokenClaims struct {
	TokenClaimsDft
	AccessTokenID string `json:"ati"`
}

// Valid checks is data in claims is valid
func (c RefreshTokenClaims) Valid() error {
	if !utils.ValidateUUID(c.GetRelatedTokenID()) {
		return bubucore.ErrTokenInvalid
	}
	return c.TokenClaimsDft.Valid()
}

// GetRelatedTokenID returns related access token ID
func (c RefreshTokenClaims) GetRelatedTokenID() string {
	return c.AccessTokenID
}
