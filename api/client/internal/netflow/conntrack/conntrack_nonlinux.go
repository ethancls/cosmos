//go:build !linux || android

package conntrack

import nftypes "github.com/ethancls/cosmos/client/internal/netflow/types"

func New(flowLogger nftypes.FlowLogger, iface nftypes.IFaceMapper) nftypes.ConnTracker {
	return nil
}
