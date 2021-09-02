package notifications

import (
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/i18n"
	"net/http"
	"net/mail"
	"strings"
)

// PlainText is a plain text notification
type PlainText struct {
	Text string `json:"text"`
}

// Filter filters notification values
func (n *PlainText) Filter() {
	n.Text = strings.TrimSpace(n.Text)
}

// Validate checks if values are OK
func (n *PlainText) Validate() error {
	n.Filter()
	if n.Text == "" {
		return bubucore.NewError(http.StatusBadRequest, "text is missing")
	}
	return nil
}

// Email notification
type Email struct {
	Recipients   []string `json:"recipients"`
	TemplateName string   `json:"template_name"`

	Subject        string                 `json:"subject,omitempty"`
	TemplateValues map[string]interface{} `json:"template_values,omitempty"`

	Language i18n.Language `json:"language,omitempty"`

	recipients []*mail.Address
}

// Filter filters notification values
func (n *Email) Filter() {
	n.Subject = strings.TrimSpace(n.Subject)
	n.TemplateName = strings.TrimSpace(n.TemplateName)
	n.Language = i18n.ParseLanguage(n.Language)
}

// Validate checks if values are OK
func (n *Email) Validate() error {
	n.Filter()
	if len(n.Recipients) == 0 {
		return bubucore.NewError(http.StatusBadRequest, "empty recipients list given")
	}
	if n.TemplateName == "" {
		return bubucore.NewError(http.StatusBadRequest, "no template name given")
	}
	_, err := n.GetRecipients()
	if err != nil {
		return bubucore.NewError(http.StatusBadRequest, "recipients list is invalid: "+err.Error())
	}
	return nil
}

// GetRecipients prepares recipients list
func (n *Email) GetRecipients() ([]*mail.Address, error) {
	if n.recipients == nil {
		var err error
		n.recipients, err = mail.ParseAddressList(strings.Join(n.Recipients, ", "))
		if err != nil {
			return nil, err
		}
	}
	return n.recipients, nil
}

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
