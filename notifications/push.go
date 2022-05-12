package notifications

import (
	"github.com/bubulearn/bubucore"
	"net/http"
	"strings"
)

// PushNotification is a FCM push notification
type PushNotification struct {
	DeviceTokens []string               `json:"device_tokens" binding:"required"`
	Title        string                 `json:"title"`
	Message      string                 `json:"message"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

// Filter request values
func (n *PushNotification) Filter() {
	n.Title = strings.TrimSpace(n.Title)
	n.Message = strings.TrimSpace(n.Message)

	dt := make([]string, 0)
	for _, t := range n.DeviceTokens {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		dt = append(dt, t)
	}

	n.DeviceTokens = dt
}

// Validate request values
func (n *PushNotification) Validate() error {
	n.Filter()
	if len(n.DeviceTokens) == 0 {
		return bubucore.NewError(http.StatusBadRequest, "empty or invalid device tokens list given")
	}
	return nil
}
