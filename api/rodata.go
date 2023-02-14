package api

import "github.com/TheDevtop/authsv/db"

// Log format
const logFormat = "api.%s: %s\n"

// Generic type
type anyType interface{}

// JSON Forms
type formReply struct {
	Error string
	Data  anyType
}

type formLogin struct {
	User   db.UserKey
	Secret db.SecretKey
}

type formQuery struct {
	User db.UserKey
	Role db.RoleKey
}

type formModifyUser struct {
	AdminUser   db.UserKey
	AdminSecret db.SecretKey
	User        db.UserKey
	Secret      db.SecretKey
	Roles       []db.RoleKey
}

type formModifyRole struct {
	AdminUser   db.UserKey
	AdminSecret db.SecretKey
	User        db.UserKey
	Role        db.RoleKey
}
