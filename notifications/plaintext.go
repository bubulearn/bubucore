package notifications

import (
	"github.com/bubulearn/bubucore"
	"net/http"
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
