package common

import (
	"sync/atomic"
	"time"

	"github.com/ethancls/kyle-server/client/firewall/manager"
	"github.com/ethancls/kyle-server/client/internal/dns"
	"github.com/ethancls/kyle-server/client/internal/peer"
	"github.com/ethancls/kyle-server/client/internal/peerstore"
	"github.com/ethancls/kyle-server/client/internal/routemanager/fakeip"
	"github.com/ethancls/kyle-server/client/internal/routemanager/iface"
	"github.com/ethancls/kyle-server/client/internal/routemanager/refcounter"
	"github.com/ethancls/kyle-server/route"
)

type HandlerParams struct {
	Route                *route.Route
	RouteRefCounter      *refcounter.RouteRefCounter
	AllowedIPsRefCounter *refcounter.AllowedIPsRefCounter
	DnsRouterInterval    time.Duration
	StatusRecorder       *peer.Status
	WgInterface          iface.WGIface
	DnsServer            dns.Server
	PeerStore            *peerstore.Store
	UseNewDNSRoute       bool
	Firewall             manager.Manager
	FakeIPManager        *fakeip.Manager
	ForwarderPort        *atomic.Uint32
}
