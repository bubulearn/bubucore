package notifications

import (
	"github.com/bubulearn/bubucore"
	"net/http"
	"net/mail"
	"strings"
)

// PlainText is a plain text notification
type PlainText struct {
	Text string `json:"text"`
}

// Email interface
type Email interface {
	GetRecipients() ([]*mail.Address, error)
}

// EmailPlain is a plain text email notification
type EmailPlain struct {
	Subject    string   `json:"subject"`
	Text       string   `json:"text"`
	Recipients []string `json:"recipients"`
}

// Filter filters notification values
func (n *EmailPlain) Filter() {
	n.Subject = strings.TrimSpace(n.Subject)
	n.Text = strings.TrimSpace(n.Text)
}

// Validate checks if values are OK
func (n *EmailPlain) Validate() error {
	n.Filter()
	if n.Subject == "" {
		return bubucore.NewError(http.StatusBadRequest, "no subject given")
	}
	if n.Text == "" {
		return bubucore.NewError(http.StatusBadRequest, "no text given")
	}
	if len(n.Recipients) == 0 {
		return bubucore.NewError(http.StatusBadRequest, "empty recipients list given")
	}
	return nil
}

// GetRecipients prepares recipients list
func (n *EmailPlain) GetRecipients() ([]*mail.Address, error) {
	return mail.ParseAddressList(strings.Join(n.Recipients, ", "))
}

// EmailWithTemplate is an email notification
type EmailWithTemplate struct {
	Subject        string      `json:"subject"`
	TemplateName   string      `json:"template_name,omitempty"`
	TemplateValues interface{} `json:"template_values,omitempty"`
	Recipients     []string    `json:"recipients"`
}

// Filter filters notification values
func (n *EmailWithTemplate) Filter() {
	n.Subject = strings.TrimSpace(n.Subject)
	n.TemplateName = strings.TrimSpace(n.TemplateName)
}

// Validate checks if values are OK
func (n *EmailWithTemplate) Validate() error {
	n.Filter()
	if n.Subject == "" {
		return bubucore.NewError(http.StatusBadRequest, "no subject given")
	}
	if n.TemplateName == "" {
		return bubucore.NewError(http.StatusBadRequest, "no template name given")
	}
	if len(n.Recipients) == 0 {
		return bubucore.NewError(http.StatusBadRequest, "empty recipients list given")
	}
	return nil
}

// GetRecipients prepares recipients list
func (n *EmailWithTemplate) GetRecipients() ([]*mail.Address, error) {
	return mail.ParseAddressList(strings.Join(n.Recipients, ", "))
}
