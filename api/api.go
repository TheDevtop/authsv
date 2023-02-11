package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TheDevtop/authsv/db"
)

// Check if user and secret matches
func authLogin(w http.ResponseWriter, r *http.Request) {
	const funcName = "authLogin"
	var form formLogin
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.Login(form.User, form.Secret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}

// List all roles of user
func authList(w http.ResponseWriter, r *http.Request) {
	const funcName = "authList"
	var form formList
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
	form.Roles = db.List(form.User)
	sendForm(w, form)
}

// Check if user matches role
func authQuery(w http.ResponseWriter, r *http.Request) {
	const funcName = "authQuery"
	var form formQuery
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.Query(form.User, form.Role); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}

// Ping server
func adminPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

// Change user secret
func adminChangeSecret(w http.ResponseWriter, r *http.Request) {
	const funcName = "adminChangeSecret"
	var form formUserMod
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.ChangeSecret(form.User, form.Secret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}

// Add new user
func adminUserAdd(w http.ResponseWriter, r *http.Request) {
	const funcName = "adminUserAdd"
	var form formUserMod
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AddUser(form.User, form.Secret, form.Roles); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}

// Delete user
func adminUserDel(w http.ResponseWriter, r *http.Request) {
	const funcName = "adminUserDel"
	var form formUserMod
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.DelUser(form.User); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}

// Add role to user
func adminRoleAdd(w http.ResponseWriter, r *http.Request) {
	const funcName = "adminRoleAdd"
	var form formRoleMod
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AddRole(form.User, form.Role); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}

// Remove role from user
func adminRoleDel(w http.ResponseWriter, r *http.Request) {
	const funcName = "adminUserDel"
	var form formRoleMod
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	} else if err := db.DelRole(form.User, form.Role); err != nil {
		log.Printf(logFmt, funcName, err)
		sendForm(w, formError{Function: funcName, Error: err.Error()})
		return
	}
}
