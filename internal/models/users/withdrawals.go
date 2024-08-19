package users

import "time"

type Withdrawals struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processedAt"`
}
