package iface

import (
	"github.com/ethancls/cosmos-server/client/iface/bind"
	"github.com/ethancls/cosmos-server/client/iface/device"
	"github.com/ethancls/cosmos-server/client/iface/netstack"
	"github.com/ethancls/cosmos-server/client/iface/wgproxy"
)

// NewWGIFace creates a new WireGuard interface for WASM (always uses netstack mode)
func NewWGIFace(opts WGIFaceOpts) (*WGIface, error) {
	relayBind := bind.NewRelayBindJS()

	wgIface := &WGIface{
		tun:            device.NewNetstackDevice(opts.IFaceName, opts.Address, opts.WGPort, opts.WGPrivKey, opts.MTU, relayBind, netstack.ListenAddr()),
		userspaceBind:  true,
		wgProxyFactory: wgproxy.NewUSPFactory(relayBind, opts.MTU),
	}

	return wgIface, nil
}
