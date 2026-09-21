package auth

import (
	"tmaster/internal/constants"
)

type Permission string

const (
	PermissionTicketCreate    Permission = "ticket:create"
	PermissionTicketReadOwn   Permission = "ticket:read_own"
	PermissionTicketReadQueue Permission = "ticket:read_queue"
	PermissionTicketAssign    Permission = "ticket:assign"
	PermissionTicketUpdate    Permission = "ticket:update"
	PermissionUserManage      Permission = "user:manage"
	PermissionTicketReadAll   Permission = "ticketread_all"
)


var RolePermission = map[constants.UserRole]map[Permission]struct{}{
	constants.RoleUser: {
		PermissionTicketCreate:  {},
		PermissionTicketReadOwn: {},
	},
	constants.RoleOperator: {
		PermissionTicketReadOwn:   {},
		PermissionTicketAssign:    {},
		PermissionTicketReadQueue: {},
		PermissionTicketUpdate:    {},
	},
	constants.RoleAdmin: {
		PermissionTicketReadAll: {},
		PermissionTicketReadQueue: {},
		PermissionUserManage: {},
	},
}


func HasPermission(role constants.UserRole, permission Permission) bool {
    permissions, exists := RolePermission[role]
    if !exists {
        return false
    }

    _, exists = permissions[permission]

    return exists
}