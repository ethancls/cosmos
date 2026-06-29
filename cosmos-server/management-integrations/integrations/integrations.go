package integrations

import (
	"context"

	"github.com/ethancls/cosmos-server/server/activity"
	"github.com/ethancls/cosmos-server/server/integrations/extra_settings"
	"github.com/ethancls/cosmos-server/server/integrations/port_forwarding"
	"github.com/ethancls/cosmos-server/server/store"
	"github.com/ethancls/cosmos-server/server/types"
)

// NewController creates a new port forwarding controller backed by the given store.
func NewController(s store.Store) port_forwarding.Controller {
	return port_forwarding.NewControllerMock()
}

// NewManager creates a new extra settings manager with the given event store.
func NewManager(eventStore activity.Store) extra_settings.Manager {
	return &extraSettingsManager{}
}

type extraSettingsManager struct {
}

func (m *extraSettingsManager) GetExtraSettings(ctx context.Context, accountID string) (*types.ExtraSettings, error) {
	return nil, nil
}

func (m *extraSettingsManager) UpdateExtraSettings(ctx context.Context, accountID, userID string, extraSettings *types.ExtraSettings) (bool, error) {
	return false, nil
}
