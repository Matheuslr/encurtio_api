package model

import "time"

type URL struct {
	ID        string
	Original  string
	ShortCode string
	CreatedAt time.Time
}
