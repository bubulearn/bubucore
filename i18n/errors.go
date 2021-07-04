package i18n

import (
	"github.com/bubulearn/bubucore"
	"net/http"
)

// ErrInvalidLang is a Language validation error
var ErrInvalidLang = bubucore.NewError(http.StatusBadRequest, "invalid language code given")
