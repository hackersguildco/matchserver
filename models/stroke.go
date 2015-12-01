package models

import "time"

//Stroke is an event point on the system
type Stroke struct {
	Location  []float32
	UserID    string
	CreatedAt time.Time
}
