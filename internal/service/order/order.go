package order

import "github.com/DenisquaP/yandex_gophermart/internal"

var accuralUrl = ""

type Order struct {
	db internal.DBStore
}

func NewOrder(db internal.DBStore, accuralSystemAddress string) *Order {
	accuralUrl = accuralSystemAddress

	return &Order{
		db: db,
	}
}
