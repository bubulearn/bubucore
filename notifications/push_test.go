package notifications_test

import (
	"github.com/bubulearn/bubucore/notifications"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
