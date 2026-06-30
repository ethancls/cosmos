package users

import (
	"github.com/ethancls/cosmos/management/server/permissions/roles"
	"github.com/ethancls/cosmos/management/server/types"
)

// Wrapped UserInfo with Role Permissions
type UserInfoWithPermissions struct {
	*types.UserInfo

	Permissions roles.Permissions
	Restricted  bool
}
