package dns

import (
	"golang.zx2c4.com/wireguard/tun/netstack"

	"github.com/ethancls/cosmos-server/client/iface/device"
	"github.com/ethancls/cosmos-server/client/iface/wgaddr"
)

// WGIface defines subset methods of interface required for manager
type WGIface interface {
	Name() string
	Address() wgaddr.Address
	IsUserspaceBind() bool
	GetFilter() device.PacketFilter
	GetDevice() *device.FilteredDevice
	GetNet() *netstack.Net
	GetInterfaceGUIDString() (string, error)
}
