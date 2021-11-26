package notifications

import (
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/utils"
	"net/http"
	"net/mail"
	"strings"
	"time"
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

// NewAppEvent creates new AppEvent notification object
func NewAppEvent(eventName string, eventData interface{}) *AppEvent {
	return &AppEvent{
		ID:        utils.GenerateUUID(),
		Name:      eventName,
		TimeEvent: time.Now().UnixNano(),
		Service:   bubucore.Opt.ServiceName,
		Data:      eventData,
	}
}

// AppEvent is a backend app event structure
type AppEvent struct {
	// Event ID
	ID string `json:"id" example:"5ce9ee33-cbc0-4603-bf40-087812c80581"`

	// Event Name
	Name string `json:"name" example:"user.new_registration"`

	// Time when event has been happened in Unix nano
	TimeEvent int64 `json:"time_event" example:"1257894000000000000"`

	// Time when event has been registered and sent to the queue
	TimeReg int64 `json:"time_reg" example:"1257894000000000000"`

	// Service is a name of a backend service where event has been happened
	Service string `json:"service" example:"calls-service"`

	// Event Data
	Data interface{} `json:"data" swaggertype:"object"`
}
