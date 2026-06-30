package account

import (
	"context"

	"github.com/ethancls/cosmos/management/server/types"
)

type RequestBuffer interface {
	GetAccountWithBackpressure(ctx context.Context, accountID string) (*types.Account, error)
}
