package db

import (
	"encoding/json"
	"fmt"
	"os"
)

// Match user and secret combination
func Login(uk UserKey, sk SecretKey) ([]RoleKey, error) {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return nil, fmt.Errorf(ErrUserNotFound, uk)
	}
	if cachedDB[uk].Secret != sk {
		return nil, fmt.Errorf(ErrLoginFailed, uk)
	}
	return cachedDB[uk].Roles, nil
}

// Match user and role combination
func Query(uk UserKey, ur RoleKey) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return fmt.Errorf(ErrUserNotFound, uk)
	}
	for _, role := range cachedDB[uk].Roles {
		if role == ur {
			return nil
		}
	}
	return fmt.Errorf(ErrRoleNotFound, uk, ur)
}

// Copy and send cachedDB
func Dump() AuthDB {
	d := make(AuthDB, len(cachedDB))
	cachedMutex.Lock()
	for k, v := range cachedDB {
		d[k] = v
	}
	cachedMutex.Unlock()
	return d
}

// Match user, secret, permission combination
// Administrative login
func AdminLogin(uk UserKey, sk SecretKey) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return fmt.Errorf(ErrUserNotFound, uk)
	} else if cachedDB[uk].Secret != sk {
		return fmt.Errorf(ErrLoginFailed, uk)
	}
	for _, role := range cachedDB[uk].Roles {
		if role == RoleAdmin {
			return nil
		}
	}
	return fmt.Errorf(ErrRoleNotFound, uk, RoleAdmin)
}

// Change secret of user
func ChangeSecret(uk UserKey, sk SecretKey) error {
	if userData, foundUser := cachedDB[uk]; !foundUser {
		return fmt.Errorf(ErrUserNotFound, uk)
	} else {
		userData.Secret = sk
		cachedMutex.Lock()
		cachedDB[uk] = userData
		cachedMutex.Unlock()
	}
	return nil
}

// Add new user
func AddUser(uk UserKey, sk SecretKey, rs []RoleKey) error {
	if _, foundUser := cachedDB[uk]; foundUser {
		return fmt.Errorf(ErrUserNotUnique, uk)
	} else if uk == "" {
		return fmt.Errorf(ErrUserNotUnique, uk)
	}
	cachedMutex.Lock()
	cachedDB[uk] = struct {
		Secret SecretKey
		Roles  []RoleKey
	}{
		Secret: sk,
		Roles:  rs,
	}
	cachedMutex.Unlock()
	return nil
}

// Delete existing user
func DelUser(uk UserKey) error {
	if _, foundUser := cachedDB[uk]; !foundUser {
		return fmt.Errorf(ErrUserNotFound, uk)
	} else if uk == "" {
		return fmt.Errorf(ErrUserNotFound, uk)
	}
	cachedMutex.Lock()
	delete(cachedDB, uk)
	cachedMutex.Unlock()
	return nil
}

// Add role to existing user
func AddRole(uk UserKey, rk RoleKey) error {
	if userData, foundUser := cachedDB[uk]; !foundUser {
		return fmt.Errorf(ErrUserNotFound, uk)
	} else if rk == "" {
		return fmt.Errorf(ErrRoleNotFound, uk, rk)
	} else {
		userData.Roles = append(userData.Roles, rk)
		cachedMutex.Lock()
		cachedDB[uk] = userData
		cachedMutex.Unlock()
	}
	return nil
}

// Delete role from existing user
func DelRole(uk UserKey, rk RoleKey) error {
	if userData, foundUser := cachedDB[uk]; !foundUser {
		return fmt.Errorf(ErrUserNotFound, uk)
	} else if rk == "" {
		return fmt.Errorf(ErrRoleNotFound, uk, rk)
	} else {
		var newRoles []RoleKey
		for _, role := range userData.Roles {
			if role == rk {
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
