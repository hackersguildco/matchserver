package models

import "time"

//Stroke is an event point on the system
type Stroke struct {
	Location  []float64 `bson:"location"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}
