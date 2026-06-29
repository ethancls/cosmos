package internal

import "github.com/ethancls/kyle-server/client/internal/stdnet"

func (e *Engine) newStdNet() (*stdnet.Net, error) {
	return stdnet.NewNetWithDiscover(e.clientCtx, e.mobileDep.IFaceDiscover, e.config.IFaceBlackList)
}
