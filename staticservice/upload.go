package staticservice

import (
	"time"
)

// Upload is an upload result struct
type Upload struct {
	ID     string `json:"id" example:"69c049e9-31cd-4103-8b83-35c391843cad"`
	UserID string `json:"user_id" example:"a44b039d-5fa2-46bf-a491-f25469b14d63"`

	Title string `json:"title" example:"filename.png"`
	Type  string `json:"type" example:"image/png"`
	Size  int64  `json:"size" example:"123000"`

	File string `json:"file" example:"69c049e9/31cd/4103/8b83/35c391843cad.png"`
	URL  string `json:"url" example:"https://static.bubulearn.dev/uploads/69c049e9/31cd/4103/8b83/35c391843cad.png"`

	TimeExpires *time.Time `json:"time_expires" bson:"time_expires" example:"2021-10-12T11:02:21.000Z"`
	TimeCreated time.Time  `json:"time_created" bson:"time_created" example:"2021-03-12T11:02:21.000Z"`
}
