package roles

import (
	"github.com/ethancls/cosmos/management/server/permissions/operations"
	"github.com/ethancls/cosmos/management/server/types"
)

var User = RolePermissions{
	Role: types.UserRoleUser,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:   false,
		operations.Create: false,
		operations.Update: false,
		operations.Delete: false,
	},
}
