package roles

import (
	"github.com/ethancls/cosmos/management/server/permissions/operations"
	"github.com/ethancls/cosmos/management/server/types"
)

var Auditor = RolePermissions{
	Role: types.UserRoleAuditor,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:   true,
		operations.Create: false,
		operations.Update: false,
		operations.Delete: false,
	},
}
