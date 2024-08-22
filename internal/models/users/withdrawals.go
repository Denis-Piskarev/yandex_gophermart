package users

import "time"

type Withdrawals struct {
	Withdrawal
	ProcessedAt time.Time `json:"processedAt"`
}

type Withdrawal struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
