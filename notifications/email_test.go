package notifications_test

import (
	"github.com/bubulearn/bubucore/notifications"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmail_Validate(t *testing.T) {
	valid := []notifications.Email{
		{
			Recipients:   []string{"test@email.com", "John Doe <john@example.com>"},
			TemplateName: "template_1",
		},
		{
			Recipients:   []string{"test@email.com", "John Doe <john@example.com>"},
			TemplateName: "template_2",
			Subject:      "Subject override",
			TemplateValues: map[string]interface{}{
				"key": "value",
			},
		},
	}

	invalid := []notifications.Email{
		{Recipients: []string{"test@email.com", "John Doe <john@example.com>"}},
		{TemplateName: "some-tpl"},
		{
			Recipients:   []string{"not an email", "John Doe <not an email too>"},
			TemplateName: "some-tpl",
		},
	}

	for _, n := range valid {
		assert.NoError(t, n.Validate())
	}

	for _, n := range invalid {
		assert.Error(t, n.Validate())
	}
}
