package db

// Types
type UserKey string

type SecretKey string

type RoleKey string

type AuthDB map[UserKey]struct {
	Secret SecretKey
	Roles  []RoleKey
}

// Constants
const RoleAdmin = "auth.admin"
const FilePerm = 0644

// Error strings
const (
	ErrUserNotFound  = "the user (%s) was not found"
	ErrRoleNotFound  = "the user (%s) does not have the role (%s)"
	ErrUserNotUnique = "the user (%s) is not unique"
	ErrLoginFailed   = "login failed for user (%s)"
)
