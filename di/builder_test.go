package di_test

import (
	"errors"
	"github.com/bubulearn/bubucore/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	b := &di.Builder{}

	err := b.Add(
		di.Def{
			Name: "test1",
			Build: func(ctn *di.Container) (interface{}, error) {
				return "testval1", nil
			},
		},
	)
	assert.NoError(t, err)

	err = b.Add(
		di.Def{
			Name: "test1",
		},
	)
	assert.Error(t, err)

	err = b.Add(
		di.Def{
			Name: "test2",
			Validate: func(ctn *di.Container) error {
				return errors.New("expected error")
			},
		},
	)
	assert.Error(t, err)

	err = b.Add(
		di.Def{
			Name: "test3",
			Build: func(ctn *di.Container) (interface{}, error) {
				return nil, errors.New("expected error")
			},
			Lazy: true,
		},
	)
	assert.NoError(t, err)

	ctn, err := b.Build()
	if !assert.NoError(t, err) {
		return
	}

	_, err = ctn.SafeGet("test3")
	assert.Error(t, err)

	v, err := ctn.SafeGet("test1")
	assert.NoError(t, err)
	assert.Equal(t, "testval1", v)

	ctn.Close()
}
