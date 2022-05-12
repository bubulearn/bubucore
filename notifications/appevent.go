package notifications

import (
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/utils"
	"time"
)

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
