package balance

import (
	"context"

	modelsBalance "github.com/DenisquaP/yandex_gophermart/internal/models/balance"
)

func (b *Balance) GetBalance(ctx context.Context, userID int) (modelsBalance.Balance, error) {
	return b.db.GetBalance(ctx, userID)
}
