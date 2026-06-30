//go:build ios

package NetBirdSDK

import (
	"github.com/ethancls/cosmos-server/util"
)

// InitializeLog initializes the log file.
func InitializeLog(logLevel string, filePath string) error {
	return util.InitLog(logLevel, filePath)
}
