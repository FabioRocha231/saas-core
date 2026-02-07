package entity

import "time"

type UserStatus string

func (u UserStatus) String() string {
	return string(u)
}

const (
	UserStatusActive  UserStatus = "active"
	UserStatusPending UserStatus = "pending"
	UserStatusBlocked UserStatus = "blocked"
)

var UserStatusMap = map[string]UserStatus{
	"active":  UserStatusActive,
	"pending": UserStatusPending,
	"blocked": UserStatusBlocked,
}

type UserRole string

func (u UserRole) String() string {
	return string(u)
}

const (
	UserRoleCostumer      UserRole = "costumer"
	UserRoleAdmin         UserRole = "admin"
	UserRoleSupport       UserRole = "support"
	UserRoleStoreOwner    UserRole = "store_owner"
	UserRoleStoreEmployee UserRole = "store_employee"
)

var UserRoleMap = map[string]UserRole{
	"costumer":       UserRoleCostumer,
	"admin":          UserRoleAdmin,
	"support":        UserRoleSupport,
	"store_owner":    UserRoleStoreOwner,
	"store_employee": UserRoleStoreEmployee,
}

type User struct {
	ID     string
	Name   string
	Cpf    string
	Email  string
	Phone  string
	Status UserStatus
	Role   UserRole

	Password string

	EmailVerifiedAt *time.Time
	PhoneVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
