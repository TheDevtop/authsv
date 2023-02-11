package db

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const (
	ErrUserNotFound  = "the specified user was not found"
	ErrRoleNotFound  = "the specified role was not found"
	ErrUserNotUnique = "the specified user is not unique"
	ErrLoginFailed   = "login failed"
	ErrUserNotAdmin  = "the specified user does not have the auth.admin role"
)

var (
	cachedDB    UserDB
	cachedMutex sync.Mutex
	cachedPath  string
)

// Match user and secret combination
func Login(uk UserKey, sk SecretKey) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	}
	if cachedDB[uk].Secret != sk {
		return errors.New(ErrLoginFailed)
	}
	return nil
}

// Return roles of user
func List(uk UserKey) RoleSet {
	return cachedDB[uk].Roles
}

// Match user and role combination
func Query(uk UserKey, ur string) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	}
	for _, role := range cachedDB[uk].Roles {
		if role == ur {
			return nil
		}
	}
	return errors.New(ErrRoleNotFound)
}

// Match user, secret, permission combination
// Administrative login
func AdminLogin(uk UserKey, sk SecretKey) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	} else if cachedDB[uk].Secret != sk {
		return errors.New(ErrLoginFailed)
	}
	for _, role := range cachedDB[uk].Roles {
		if role == RoleAdmin {
			return nil
		}
	}
	return errors.New(ErrUserNotAdmin)
}

// Change secret of user
func ChangeSecret(uk UserKey, sk SecretKey) error {
	if userData, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	} else {
		userData.Secret = sk
		cachedMutex.Lock()
		cachedDB[uk] = userData
		cachedMutex.Unlock()
	}
	return nil
}

// Add new user
func AddUser(uk UserKey, sk SecretKey, roles RoleSet) error {
	if _, foundUser := cachedDB[uk]; foundUser {
		return errors.New(ErrUserNotUnique)
	} else if uk == "" {
		return errors.New(ErrUserNotUnique)
	}
	cachedMutex.Lock()
	cachedDB[uk] = struct {
		Secret SecretKey
		Roles  RoleSet
	}{
		Secret: sk,
		Roles:  roles,
	}
	cachedMutex.Unlock()
	return nil
}

// Delete existing user
func DelUser(uk UserKey) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	} else if uk == "" {
		return errors.New(ErrUserNotFound)
	}
	cachedMutex.Lock()
	delete(cachedDB, uk)
	cachedMutex.Unlock()
	return nil
}

// Add role to existing user
func AddRole(uk UserKey, ur string) error {
	if userData, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	} else if ur == "" {
		return errors.New(ErrRoleNotFound)
	} else {
		userData.Roles = append(userData.Roles, ur)
		cachedMutex.Lock()
		cachedDB[uk] = userData
		cachedMutex.Unlock()
	}
	return nil
}

// Delete role from existing user
func DelRole(uk UserKey, ur string) error {
	if userData, foundUser := cachedDB[uk]; !foundUser {
		return errors.New(ErrUserNotFound)
	} else if ur == "" {
		return errors.New(ErrRoleNotFound)
	} else {
		var newRoles RoleSet
		for _, role := range userData.Roles {
			if role == ur {
				continue
			}
			newRoles = append(newRoles, role)
		}
		userData.Roles = newRoles
		cachedMutex.Lock()
		cachedDB[uk] = userData
		cachedMutex.Unlock()
	}
	return nil
}

// Load cached database from filesystem
func LoadCache(path string) error {
	if buf, err := os.ReadFile(path); err != nil {
		return err
	} else if err := json.Unmarshal(buf, &cachedDB); err != nil {
		cachedDB = nil
		return err
	}
	cachedPath = path
	return nil
}

// Store cached database to filesystem
func StoreCache() error {
	if buf, err := json.Marshal(cachedDB); err != nil {
		return err
	} else if err := os.WriteFile(cachedPath, buf, FilePerm); err != nil {
		return err
	}
	return nil
}
