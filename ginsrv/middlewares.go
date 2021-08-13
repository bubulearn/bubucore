package ginsrv

import (
	"bytes"
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/di"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/tokens"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// KeyAccessClaims is a param name for current access claims
	KeyAccessClaims = "BubuAccessClaims"

	// KeyI18nLang is a context key for language
	KeyI18nLang = "BubuI18nLang"

	// KeyDIContainer is a context key for a DI di.Container
	KeyDIContainer = "BubuDIContainer"
)

// M returns Middlewares instance
func M() *Middlewares {
	if middlewares == nil {
		middlewares = &Middlewares{}
	}
	return middlewares
}

// middlewares is a current Middlewares instance
var middlewares *Middlewares

// Middlewares contains middlewares functions
type Middlewares struct {
}

// SetDIContainer is a middleware to set app.App instance to the context
func (m *Middlewares) SetDIContainer(ctn *di.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDIContainer, ctn)
	}
}

// JWTAccess is a authorization by the Access token.
// Sets parsed claims to KeyAccessClaims param.
func (m *Middlewares) JWTAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)

		var err error
		sign, err := ctx.ExtractBearerToken()
		if err != nil {
			e := bubucore.NewError(http.StatusUnauthorized, err.Error())
			ctx.Err(e)
			ctx.Abort()
			return
		}

		claims, err := tokens.ParseAccessToken(sign)
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		ctx.Set(KeyAccessClaims, claims)

		err = claims.Valid()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}

		if claims.Language != "" {
			ctx.SetI18nLang(claims.Language)
		}

		ctx.Next()
	}
}

// RequireRole validates user role is equal to specified
func (m *Middlewares) RequireRole(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		if claims.Role != role {
			ctx.Err(bubucore.ErrRoleNotAllowed)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// RequireRoleIn validates if user role is one of specified
func (m *Middlewares) RequireRoleIn(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		for _, role := range roles {
			if claims.Role == role {
				ctx.Next()
				return
			}
		}
		ctx.Err(bubucore.ErrRoleNotAllowed)
		ctx.Abort()
	}
}

// RequireRoleHigher validates if user role is equal or higher of specified
func (m *Middlewares) RequireRoleHigher(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		if claims.Role < role {
			ctx.Err(bubucore.ErrRoleNotAllowed)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// InitI18nLang reads language from Accept-Language header and sets it to the context
func (m *Middlewares) InitI18nLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		accept := ctx.GetHeader("Accept-Language")
		lang := i18n.ParseAcceptLanguageString(accept)
		ctx.SetI18nLang(lang)
	}
}

// LogBody is a middleware to write response body to the log
func (m *Middlewares) LogBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writer := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = writer
		ctx.Next()
		logger := log.WithFields(log.Fields{
			bubucore.LogFieldType:   bubucore.LogTypeHTTPIO,
			bubucore.LogFieldStatus: ctx.Writer.Status(),
			bubucore.LogFieldPath:   ctx.FullPath(),
			bubucore.LogFieldMethod: ctx.Request.Method,
		})
		if ctx.Writer.Status() >= 500 {
			logger.Error(writer.body.String())
		} else if ctx.Writer.Status() >= 400 {
			logger.Warn(writer.body.String())
		}
	}
}

// LogFormatter formats gin log record as JSON string
func (m *Middlewares) LogFormatter(param gin.LogFormatterParams) string {
	data := map[string]string{
		"@timestamp":            param.TimeStamp.Format(time.RFC3339),
		"ip":                    param.ClientIP,
		bubucore.LogFieldMethod: param.Method,
		bubucore.LogFieldPath:   param.Path,
		"proto":                 param.Request.Proto,
		bubucore.LogFieldStatus: strconv.FormatInt(int64(param.StatusCode), 10),
		"latency":               strconv.FormatFloat(param.Latency.Seconds(), 'f', 8, 64),
		"latency_fmt":           param.Latency.String(),
		"agent":                 param.Request.UserAgent(),
		"error":                 param.ErrorMessage,
		"request_body":          "-",
		"response_body_size":    strconv.FormatInt(int64(param.BodySize), 10),
	}

	if param.StatusCode >= 400 && param.Request.Body != nil {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, param.Request.Body)
		data["request_body"] = buf.String()
		defer func() {
			_ = param.Request.Body.Close()
		}()
	}

	data[bubucore.LogFieldType] = bubucore.LogTypeHTTPSrv
	data[bubucore.LogFieldService] = bubucore.Opt.ServiceName
	data[bubucore.LogFieldHostname] = bubucore.Opt.GetHostname()
	data[bubucore.LogFieldAPIVersion] = bubucore.Opt.APIVersion

	if param.StatusCode >= 500 {
		data[bubucore.LogFieldLevel] = "error"
	} else if param.StatusCode >= 400 {
		data[bubucore.LogFieldLevel] = "warning"
	} else {
		data[bubucore.LogFieldLevel] = "info"
	}

	// not using marshalling for speed-up
	values := make([]string, len(data))
	i := 0
	for key, val := range data {
		val = strings.ReplaceAll(val, `"`, `\"`)
		val = strings.TrimSpace(val)
		values[i] = `"` + key + `":"` + val + `"`
		i++
	}

	return "{" + strings.Join(values, ",") + "}\n"
}

// bodyLogWriter is a writer to write response body to the log
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write body
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
