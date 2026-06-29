package config

import (
	"github.com/ethancls/cosmos-server/shared/management/proto"
	"github.com/ethancls/cosmos-server/server/types"
)

// ExtendCosmosConfig extends the Cosmos config with integration-specific settings.
func ExtendCosmosConfig(peerID string, peerGroups []string, nbConfig *proto.NetbirdConfig, extraSettings *types.ExtraSettings) *proto.NetbirdConfig {
	// For now, return the config as-is. Integration extensions can be added here.
	return nbConfig
}
