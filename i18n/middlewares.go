package i18n

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// GinKeyLang is a context key for language
const GinKeyLang = "BubuLang"

// MiddlewareInitLang reads language from Accept-Language header and sets it to the context
func MiddlewareInitLang() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accept := ctx.GetHeader("Accept-Language")
		lang := parseAcceptLanguageString(accept)
		GinSetLang(ctx, lang)
	}
}

// GinGetLang gets language from gin context
func GinGetLang(ctx *gin.Context) Language {
	v, ok := ctx.Get(GinKeyLang)
	if !ok {
		return Source.DefaultLang
	}
	lang, ok := v.(Language)
	if !ok || lang == "" {
		return Source.DefaultLang
	}
	return lang
}

// GinSetLang set lang for gix context
func GinSetLang(ctx *gin.Context, lang Language) {
	ctx.Set(GinKeyLang, lang)
}

// parseAcceptLanguageString does simple parsing of the first language
func parseAcceptLanguageString(accept string) Language {
	parts := strings.Split(accept, ", ")
	p := strings.TrimSpace(parts[0])
	if p == "" {
		return Source.DefaultLang
	}

	parts = strings.Split(accept, "-")
	return Language(strings.ToLower(parts[0]))
}
