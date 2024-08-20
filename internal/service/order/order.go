package order

import "github.com/DenisquaP/yandex_gophermart/internal"

var accuralURL = ""

type Order struct {
	db internal.DBStore
}

func NewOrder(db internal.DBStore, accuralSystemAddress string) *Order {
	accuralURL = accuralSystemAddress

	return &Order{
		db: db,
	}
}
