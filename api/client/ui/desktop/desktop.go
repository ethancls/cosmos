package desktop

import "github.com/ethancls/cosmos/version"

// GetUIUserAgent returns the Desktop ui user agent
func GetUIUserAgent() string {
	return "netbird-desktop-ui/" + version.NetbirdVersion()
}
