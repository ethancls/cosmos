package integrations

import (
	"context"

	"github.com/ethancls/cosmos/management/server/activity"
	"github.com/ethancls/cosmos/management/server/types"

	"github.com/ethancls/cosmos/management/server/integrations/extra_settings"
)

type ManagerImpl struct {
}

func NewManager(eventStore activity.Store) extra_settings.Manager {
	return &ManagerImpl{}
}

func (m *ManagerImpl) GetExtraSettings(ctx context.Context, accountID string) (*types.ExtraSettings, error) {
	return &types.ExtraSettings{}, nil
}

func (m *ManagerImpl) UpdateExtraSettings(ctx context.Context, accountID, userID string, accountExtraSettings *types.ExtraSettings) (bool, error) {
	return false, nil
}
