package notifications

import (
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/utils"
	"net/http"
	"strings"
)

// SMS is an SMS notification struct
type SMS struct {
	To   string `json:"to" example:"+79001234567"`
	Text string `json:"text" example:"Hello, John!"`
}

// Filter filters request values
func (r *SMS) Filter() {
	r.To = strings.TrimSpace(r.To)
	r.Text = strings.TrimSpace(r.Text)
}

// Validate validates request values
func (r *SMS) Validate() error {
	r.Filter()

	if r.Text == "" {
		return bubucore.NewError(http.StatusPreconditionFailed, "no sms text in `text` field given")
	}

	if !utils.ValidatePhone(r.To) {
		return bubucore.NewError(http.StatusPreconditionFailed, "invalid sms phone number in `to` field")
	}

	r.To = utils.NormalizePhone(r.To)

	return nil
}
