package roles

import (
	"github.com/ethancls/cosmos/management/server/permissions/modules"
	"github.com/ethancls/cosmos/management/server/permissions/operations"
	"github.com/ethancls/cosmos/management/server/types"
)

type RolePermissions struct {
	Role         types.UserRole
	Permissions  Permissions
	AutoAllowNew map[operations.Operation]bool
}

type Permissions map[modules.Module]map[operations.Operation]bool

var RolesMap = map[types.UserRole]RolePermissions{
	types.UserRoleOwner:        Owner,
	types.UserRoleAdmin:        Admin,
	types.UserRoleUser:         User,
	types.UserRoleAuditor:      Auditor,
	types.UserRoleNetworkAdmin: NetworkAdmin,
}
