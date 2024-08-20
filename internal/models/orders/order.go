package orders

import "time"

// Order - struct of order
type Order struct {
	Number     int       `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploadedAt,omitempty"`
}

type OrderAccrual struct {
	Order   int     `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}
