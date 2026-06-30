package config

import (
	"github.com/ethancls/cosmos/management/server/types"
	"github.com/ethancls/cosmos/shared/management/proto"
)

func ExtendNetBirdConfig(peerID string, peerGroups []string, config *proto.NetbirdConfig, extraSettings *types.ExtraSettings) *proto.NetbirdConfig {
	return config
}
