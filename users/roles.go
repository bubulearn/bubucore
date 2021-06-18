package models

import (
	"github.com/bubulearn/bubucore"
	"net/http"
)

// Users' roles
const (
	RoleStudent = 1
	RoleTeacher = 500
	RoleAdmin   = 1000
)

var rolesAvailable = []int{
	RoleStudent,
	RoleTeacher,
	RoleAdmin,
}

// ValidateRole validates role
func ValidateRole(role int) error {
	roleValid := false
	for _, r := range rolesAvailable {
		if r == role {
			roleValid = true
			break
		}
	}
	if !roleValid {
		return bubucore.NewError(http.StatusBadRequest, "user role is not valid")
	}
	return nil
}
