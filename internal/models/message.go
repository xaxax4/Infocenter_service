package models

import "time"

type Message struct {
	ID      int
	Content string
	Time    time.Time
}
