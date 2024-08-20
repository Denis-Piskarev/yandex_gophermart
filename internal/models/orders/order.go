package orders

import "time"

// Order - struct of order
type Order struct {
	Number     int       `json:"number"`
	Status     string    `json:"status"`
	Accural    int       `json:"accural,omitempty"`
	UploadedAt time.Time `json:"uploadedAt,omitempty"`
}
