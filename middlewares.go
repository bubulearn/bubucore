package bubucore

import (
	"bytes"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// KeyAccessClaims is a param name for current access claims
const KeyAccessClaims = "BubuAccessClaims"

// MiddlewareJWTAccess is a authorization by the Access token.
// Sets parsed claims to KeyAccessClaims param.
func MiddlewareJWTAccess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		sign, err := utils.ExtractBearerToken(ctx.Request)
		if err != nil {
			e := NewError(http.StatusUnauthorized, err.Error())
			ErrorResponseE(ctx, e, http.StatusForbidden)
			ctx.Abort()
			return
		}

		claims, err := ParseAccessToken(sign)
		if err != nil {
			ErrorResponseE(ctx, err, http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		ctx.Set(KeyAccessClaims, claims)

		err = claims.Valid()
		if err != nil {
			ErrorResponseE(ctx, err, http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		if claims.Language != "" {
			i18n.GinSetLang(ctx, claims.Language)
		}

		ctx.Next()
	}
}

// MiddlewareValidateRole validates user role from claims
func MiddlewareValidateRole(role int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := GinGetAccessClaims(ctx)
		if err != nil {
			ErrorResponseE(ctx, err, http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		if claims.Role != role {
			ErrorResponseS(ctx, "unexpected user role", http.StatusForbidden)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// MiddlewareLogBody is a middleware to write response body to the log
func MiddlewareLogBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writer := &ginBodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = writer
		ctx.Next()
		logger := log.WithField(LogFieldType, LogTypeHTTPIO)
		if ctx.Writer.Status() >= 500 {
			logger.Error(writer.body.String())
		} else if ctx.Writer.Status() >= 400 {
			logger.Warn(writer.body.String())
		}
	}
}

// ginBodyLogWriter is a writer to write response body to the log
type ginBodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write body
func (w ginBodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinGetAccessClaims returns AccessTokenClaims from the current gin context
func GinGetAccessClaims(ctx *gin.Context) (*AccessTokenClaims, error) {
	c, ok := ctx.Get(KeyAccessClaims)
	if !ok {
		log.Warn("no claims initialized")
		return nil, ErrTokenInvalid
	}
	claims, ok := c.(*AccessTokenClaims)
	if !ok {
		log.Warn("unexpected claims type")
		return nil, ErrTokenInvalid
	}
	return claims, nil
}
