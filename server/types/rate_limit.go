package types

import "time"

type RateLimit struct {
	Tokens   uint64        `json:"tokens"`
	Interval time.Duration `json:"interval"`
}
