package users

import "time"

type Withdrawals struct {
	Order       int       `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processedAt"`
}
