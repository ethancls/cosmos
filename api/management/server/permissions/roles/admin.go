package roles

import (
	"github.com/ethancls/cosmos/management/server/permissions/modules"
	"github.com/ethancls/cosmos/management/server/permissions/operations"
	"github.com/ethancls/cosmos/management/server/types"
)

var Admin = RolePermissions{
	Role: types.UserRoleAdmin,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:   true,
		operations.Create: true,
		operations.Update: true,
		operations.Delete: true,
	},
	Permissions: Permissions{
		modules.Accounts: {
			operations.Read:   true,
			operations.Create: false,
			operations.Update: false,
			operations.Delete: false,
		},
	},
}
