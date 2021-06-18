package utils

import (
	"errors"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

var regxPhone *regexp.Regexp

// ExtractBearerToken extracts token from the 'Authorization: Bearer <token>' header
func ExtractBearerToken(r *http.Request) (string, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", errors.New("no Authorization header provided")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("unsupported Authorization header format")
	}

	return parts[1], nil
}

// GenerateUUID creates new id string
func GenerateUUID() string {
	return uuid.New().String()
}

// ValidateUUID validates id
func ValidateUUID(id string) bool {
	if id == "" {
		return false
	}
	_, err := uuid.Parse(id)
	return err == nil
}

// FilterEmail value
func FilterEmail(email string) string {
	email = strings.TrimSpace(email)
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ""
	}
	return addr.Address
}

// ValidateEmail validates email value
func ValidateEmail(email string) bool {
	email = FilterEmail(email)
	return email != ""
}

// ValidatePhone validates phone number
func ValidatePhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false
	}
	if regxPhone == nil {
		regxPhone = regexp.MustCompile(`^[+0-9\s\-()]{5,}$`)
	}
	return regxPhone.MatchString(phone)
}

// JSONConvert converts anything to anything through JSON marshalling
func JSONConvert(source interface{}, target interface{}) error {
	json, err := jsoniter.Marshal(source)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(json, target)
}
