//go:build !android

package internal

import (
	"github.com/ethancls/cosmos/client/internal/stdnet"
)

func (e *Engine) newStdNet() (*stdnet.Net, error) {
	return stdnet.NewNet(e.clientCtx, e.config.IFaceBlackList)
}
