package notifications_test

import (
	"github.com/bubulearn/bubucore/notifications"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSMS_Validate(t *testing.T) {
	{
		t.Log("testing valid requests")

		requests := []*notifications.SMS{
			{
				To:   "+7 (900) 123-45-67",
				Text: "What keeps you up nights, Mr. Dillinger?",
			},
			{
				To:   "8900456-78-90",
				Text: "Coffee.",
			},
		}

		for i, req := range requests {
			assert.NoError(t, req.Validate(), i)
		}
	}

	{
		t.Log("testing invalid requests")

		requests := []*notifications.SMS{
			{
				To:   "What is it?",
				Text: "What keeps you up nights, Mr. Dillinger?",
			},
			{
				To:   "8900456-78-90",
				Text: "",
			},
		}

		for i, req := range requests {
			assert.Error(t, req.Validate(), i)
		}
	}
}
