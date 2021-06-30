package i18n

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

// T is a Translations.GetText alias
func T(textOrKey string, lang Language, tplData interface{}) string {
	return Source.Translations.GetText(textOrKey, lang, tplData)
}
