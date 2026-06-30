//go:build !android

package server

import (
	"github.com/ethancls/cosmos/client/firewall/iptables"
	"github.com/ethancls/cosmos/client/firewall/nftables"
	"github.com/ethancls/cosmos/client/internal/dns"
	"github.com/ethancls/cosmos/client/internal/routemanager/systemops"
	"github.com/ethancls/cosmos/client/internal/statemanager"
	"github.com/ethancls/cosmos/client/ssh/config"
)

// registerStates registers all states that need crash recovery cleanup.
func registerStates(mgr *statemanager.Manager) {
	mgr.RegisterState(&dns.ShutdownState{})
	mgr.RegisterState(&systemops.ShutdownState{})
	mgr.RegisterState(&nftables.ShutdownState{})
	mgr.RegisterState(&iptables.ShutdownState{})
	mgr.RegisterState(&config.ShutdownState{})
}
