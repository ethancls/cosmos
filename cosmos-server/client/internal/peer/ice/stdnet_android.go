package ice

import (
	"context"

	"github.com/ethancls/cosmos-server/client/internal/stdnet"
)

func newStdNet(ctx context.Context, iFaceDiscover stdnet.ExternalIFaceDiscover, ifaceBlacklist []string) (*stdnet.Net, error) {
	return stdnet.NewNetWithDiscover(ctx, iFaceDiscover, ifaceBlacklist)
}
