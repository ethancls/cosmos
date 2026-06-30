//go:build !linux && !ios && !android && !js

package iface

import (
	"github.com/ethancls/cosmos-server/client/iface/bind"
	"github.com/ethancls/cosmos-server/client/iface/device"
	"github.com/ethancls/cosmos-server/client/iface/netstack"
	"github.com/ethancls/cosmos-server/client/iface/wgproxy"
)

// NewWGIFace Creates a new WireGuard interface instance
func NewWGIFace(opts WGIFaceOpts) (*WGIface, error) {
	iceBind := bind.NewICEBind(opts.TransportNet, opts.Address, opts.MTU)

	var tun WGTunDevice
	if netstack.IsEnabled() {
		tun = device.NewNetstackDevice(opts.IFaceName, opts.Address, opts.WGPort, opts.WGPrivKey, opts.MTU, iceBind, netstack.ListenAddr())
	} else {
		tun = device.NewTunDevice(opts.IFaceName, opts.Address, opts.WGPort, opts.WGPrivKey, opts.MTU, iceBind)
	}

	return &WGIface{
		userspaceBind:  true,
		tun:            tun,
		wgProxyFactory: wgproxy.NewUSPFactory(iceBind, opts.MTU),
	}, nil
}
