package models

import "time"

type Client struct {
	Chat      chan Message
	CreatedAt time.Time
}
