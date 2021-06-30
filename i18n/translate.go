package i18n

// Translatable are structures which can translate their values
type Translatable interface {
	Translate()
}

// TranslateSlice calls Translate on each Translatable item in slice
func TranslateSlice(items []Translatable) {
	for _, item := range items {
		item.Translate()
	}
}

// TranslateStringMap calls Translate on each Translatable item in string map
func TranslateStringMap(items map[string]Translatable) {
	for _, item := range items {
		item.Translate()
	}
}
