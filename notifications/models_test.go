package notifications_test

import (
	"github.com/bubulearn/bubucore/notifications"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlainText_Filter(t *testing.T) {
	values := map[string]string{
		"   ":     "",
		" test ":  "test",
		"Test  1": "Test  1",
	}

	for in, expected := range values {
		n := &notifications.PlainText{
			Text: in,
		}
		n.Filter()
		assert.Equal(t, expected, n.Text)
	}
}

func TestPlainText_Validate(t *testing.T) {
	valid := []string{
		"Test 1",
		"<pre>test 2</pre>",
	}
	invalid := []string{
		"",
		"    ",
		"\n\n",
	}

	for _, v := range valid {
		n := &notifications.PlainText{
			Text: v,
		}
		assert.NoError(t, n.Validate())
	}

	for _, v := range invalid {
		n := &notifications.PlainText{
			Text: v,
		}
		assert.Error(t, n.Validate())
	}
}

func TestEmail_Validate(t *testing.T) {
	valid := []notifications.Email{
		{
			Recipients:   []string{"test@email.com", "John Doe <john@example.com>"},
			TemplateName: notifications.TplResetPassLink,
		},
		{
			Recipients:   []string{"test@email.com", "John Doe <john@example.com>"},
			TemplateName: notifications.TplResetPassLink,
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

func TestPushNotification_Validate(t *testing.T) {
	{
		t.Log("Testing valid requests")
		valid := []*notifications.PushNotification{
			{
				DeviceTokens: []string{"test1", " test2 "},
				Title:        "Hello",
				Message:      "World",
				Data: map[string]interface{}{
					"key": "value",
				},
			},
			{
				DeviceTokens: []string{"test1", " test2 "},
				Title:        "Hello",
				Message:      "World",
			},
		}
		for i, n := range valid {
			assert.NoError(t, n.Validate(), "Row #", i)
		}
	}

	{
		t.Log("Testing invalid requests")
		invalid := []*notifications.PushNotification{
			{},
			{
				DeviceTokens: []string{"  ", ""},
			},
			{
				DeviceTokens: []string{},
			},
		}
		for i, n := range invalid {
			assert.Error(t, n.Validate(), "Row #", i)
		}
	}
}
