package constants

type UserRole string

const (
	RoleUser     UserRole = "user"
	RoleOperator UserRole = "operator"
	RoleAdmin    UserRole = "admin"
)

type TicketStatus string

const (
	StatusNew     TicketStatus = "new"
	StatusPending TicketStatus = "pending"
	StatusClosed  TicketStatus = "closed"
)

type contextKey string

const (
	RequestIDKey  contextKey = "request_id"
	UserIDKey     contextKey = "userID"
	UserRoleKey   contextKey = "userRole"
	RequestLogger contextKey = "requestLogger"
)
