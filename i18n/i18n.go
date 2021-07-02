package i18n

import (
	"strings"
)

// Available languages
const (
	LangRu = Language("ru")
	LangEn = Language("en")
)

// Source is a current texts source
var Source = &TextsSource{
	DefaultLang:  LangEn,
	Translations: make(Translations),
}

// Language is a language code
type Language string

// Filter value
func (l *Language) Filter() {
	*l = Language(strings.TrimSpace(string(*l)))
}

// Validate if value is ok
func (l Language) Validate() error {
	return nil
}

// T is a Translations.GetText alias
func T(textOrKey string, lang Language, tplData interface{}) string {
	return Source.Translations.GetText(textOrKey, lang, tplData)
}
