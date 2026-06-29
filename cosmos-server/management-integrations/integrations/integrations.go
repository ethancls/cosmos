package integrations

import (
	"github.com/ethancls/cosmos-server/server/activity"
	"github.com/ethancls/cosmos-server/server/integrations/extra_settings"
	"github.com/ethancls/cosmos-server/server/integrations/port_forwarding"
	"github.com/ethancls/cosmos-server/server/store"
)

// NewController creates a new port forwarding controller backed by the given store.
func NewController(s store.Store) port_forwarding.Controller {
	return port_forwarding.NewControllerMock()
}

// NewManager creates a new extra settings manager with the given event store.
func NewManager(eventStore *activity.InMemoryEventStore) extra_settings.Manager {
	return &extraSettingsManager{eventStore: eventStore}
}

type extraSettingsManager struct {
	eventStore *activity.InMemoryEventStore
}

func (m *extraSettingsManager) GetExtraSettings(ctx interface{}, accountID string) (*extra_settings.ExtraSettings, error) {
	return nil, nil
}

func (m *extraSettingsManager) UpdateExtraSettings(ctx interface{}, accountID, userID string, extraSettings *extra_settings.ExtraSettings) (bool, error) {
	return false, nil
}
