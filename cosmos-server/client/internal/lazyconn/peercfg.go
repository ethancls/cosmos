package lazyconn

import (
	"net/netip"

	log "github.com/sirupsen/logrus"

	"github.com/ethancls/kyle-server/client/internal/peer/id"
)

type PeerConfig struct {
	PublicKey  string
	AllowedIPs []netip.Prefix
	PeerConnID id.ConnID
	Log        *log.Entry
}
