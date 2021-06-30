package i18n

// Translatable are structures which can translate their values
type Translatable interface {
	Translate(lang Language)
}

// TranslateSlice calls Translate on each Translatable item in slice
func TranslateSlice(items []Translatable, lang Language) {
	for _, item := range items {
		item.Translate(lang)
	}
}

// TranslateStringMap calls Translate on each Translatable item in string map
func TranslateStringMap(items map[string]Translatable, lang Language) {
	for _, item := range items {
		item.Translate(lang)
	}
}
