package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateRole(t *testing.T) {
	valid := rolesAvailable
	invalid := []int{
		42,
		-1,
		0,
	}

	for _, r := range valid {
		err := ValidateRole(r)
		assert.NoError(t, err)
	}

	for _, r := range invalid {
		err := ValidateRole(r)
		assert.Error(t, err)
	}
}
