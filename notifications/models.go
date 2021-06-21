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

// Email notification
type Email struct {
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`

	Text string `json:"text,omitempty"`

	TemplateName   string      `json:"template_name,omitempty"`
	TemplateValues interface{} `json:"template_values,omitempty"`
}

// Filter filters notification values
func (n *Email) Filter() {
	n.Subject = strings.TrimSpace(n.Subject)
	n.Text = strings.TrimSpace(n.Text)
	n.TemplateName = strings.TrimSpace(n.TemplateName)
}

// Validate checks if values are OK
func (n *Email) Validate() error {
	n.Filter()
	if n.Subject == "" {
		return bubucore.NewError(http.StatusBadRequest, "no subject given")
	}
	if len(n.Recipients) == 0 {
		return bubucore.NewError(http.StatusBadRequest, "empty recipients list given")
	}
	if n.Text == "" && n.TemplateName == "" {
		return bubucore.NewError(http.StatusBadRequest, "no text or template name given")
	}
	return nil
}

// GetRecipients prepares recipients list
func (n *Email) GetRecipients() ([]*mail.Address, error) {
	return mail.ParseAddressList(strings.Join(n.Recipients, ", "))
}
