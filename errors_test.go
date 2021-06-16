package bubucore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(400, "test")
	assert.Error(t, err)
	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "test", err.Message)
}

func TestError_Error(t *testing.T) {
	errs := map[string]*Error{
		"test1": NewError(1, "test1"),
		"test2": NewError(2, "test2"),
	}

	for text, err := range errs {
		assert.Equal(t, text, err.Error())
	}
}
