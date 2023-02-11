package db

const RoleAdmin = "auth.admin"

const FilePerm = 0644

type UserKey string

type SecretKey string

type RoleSet []string

type UserDB map[UserKey]struct {
	Secret SecretKey
	Roles  RoleSet
}
