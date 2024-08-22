package order

import "github.com/DenisquaP/yandex_gophermart/internal"

var accrualURL = ""

type Order struct {
	db internal.DBStore
}

func NewOrder(db internal.DBStore, accrualSystemAddress string) *Order {
	accrualURL = accrualSystemAddress

	return &Order{
		db: db,
	}
}
