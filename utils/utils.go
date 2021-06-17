package utils

import (
	"github.com/bubulearn/bubucore"
	"net/http"
	"strings"
)

// ExtractBearerToken extracts token from the 'Authorization: Bearer <token>' header
func ExtractBearerToken(r *http.Request) (string, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", bubucore.NewError(http.StatusUnauthorized, "no Authorization header provided")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", bubucore.NewError(http.StatusUnauthorized, "unsupported Authorization header format")
	}

	return parts[1], nil
}
