package models

import "time"

type Nonce struct {
	T        time.Time
	S        string
	Endpoint string
}
