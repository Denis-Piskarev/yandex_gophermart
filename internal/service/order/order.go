package order

import "github.com/DenisquaP/yandex_gophermart/internal"

type Order struct {
	db internal.DBStore
}

func NewOrder(db internal.DBStore) *Order {
	return &Order{
		db: db,
	}
}
