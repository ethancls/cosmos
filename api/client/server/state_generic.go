//go:build !linux || android

package server

import (
	"github.com/ethancls/cosmos/client/internal/dns"
	"github.com/ethancls/cosmos/client/internal/routemanager/systemops"
	"github.com/ethancls/cosmos/client/internal/statemanager"
	"github.com/ethancls/cosmos/client/ssh/config"
)

// registerStates registers all states that need crash recovery cleanup.
func registerStates(mgr *statemanager.Manager) {
	mgr.RegisterState(&dns.ShutdownState{})
	mgr.RegisterState(&systemops.ShutdownState{})
	mgr.RegisterState(&config.ShutdownState{})
}
