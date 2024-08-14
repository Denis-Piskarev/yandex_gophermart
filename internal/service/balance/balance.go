package balance

import "github.com/DenisquaP/yandex_gophermart/internal"

type Balance struct {
	db internal.DBStore
}

func NewBalance(db internal.DBStore) *Balance {
	return &Balance{
		db: db,
	}
}
