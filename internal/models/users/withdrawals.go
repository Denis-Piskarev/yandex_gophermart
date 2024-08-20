package users

import "time"

type Withdrawals struct {
	Withdrawal
	ProcessedAt time.Time `json:"processedAt"`
}

type Withdrawal struct {
	Order float64 `json:"order"`
	Sum   int     `json:"sum"`
}
