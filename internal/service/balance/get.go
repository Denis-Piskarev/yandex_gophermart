package balance

import (
	"context"

	modelsBalance "github.com/DenisquaP/yandex_gophermart/internal/models/balance"
)

func (b *Balance) GetBalance(ctx context.Context, userId int) (modelsBalance.Balance, error) {
	return b.db.GetBalance(ctx, userId)
}
