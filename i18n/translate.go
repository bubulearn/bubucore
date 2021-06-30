package i18n

// Translatable are structures which can translate their values
type Translatable interface {
	Translate(lang Language)
}
