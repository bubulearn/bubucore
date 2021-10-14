package bubucore_test

import (
	"github.com/bubulearn/bubucore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := bubucore.NewError(400, "test")
	assert.Error(t, err)
	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "test", err.Message)

	err = bubucore.NewError(http.StatusNotFound)
	assert.Error(t, err)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "Not Found", err.Message)
}

func TestError_Error(t *testing.T) {
	errs := map[string]*bubucore.Error{
		"test1": bubucore.NewError(1, "test1"),
		"test2": bubucore.NewError(2, "test2"),
		"test3": bubucore.NewError(2, "test", 3),
	}

	for text, err := range errs {
		assert.Equal(t, text, err.Error())
	}
}
