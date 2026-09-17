package auth

type Permission string

const (
	PermissionTicketCreate    Permission = "ticket:create"
	PermissionTicketReadOwn   Permission = "ticket:read_own"
	PermissionTicketReadQueue Permission = "ticket:read_queue"
	PermissionTicketAssign    Permission = "ticket:assign"
	PermissionTicketUpdate    Permission = "ticket:update"
	PermissionUserManage      Permission = "user:manage"
)

type UserRole string

const (
	RoleUser     UserRole = "user"
	RoleOperator UserRole = "operator"
	RoleAdmin    UserRole = "admin"
)

var RolePermission = map[UserRole]map[Permission]struct{} {
    RoleUser: {

    },
    RoleOperator: {

    },
    RoleAdmin: {

    },
}