package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TheDevtop/authsv/db"
)

type formError struct {
	Function string
	Error    string
}

type formLogin struct {
	User   db.UserKey
	Secret db.SecretKey
}

type formQuery struct {
	User db.UserKey
	Role string
}

type formList struct {
	User  db.UserKey
	Roles db.RoleSet
}

type formUserMod struct {
	AdminUser   db.UserKey
	AdminSecret db.SecretKey
	User        db.UserKey
	Secret      db.SecretKey
	Roles       db.RoleSet
}

type formRoleMod struct {
	AdminUser   db.UserKey
	AdminSecret db.SecretKey
	User        db.UserKey
	Role        string
}

// Receive form from client, via pointer
func receiveForm(r *http.Request, formPtr interface{}) error {
	if buf, err := io.ReadAll(r.Body); err != nil {
		return err
	} else if err := json.Unmarshal(buf, formPtr); err != nil {
		return err
	}
	return nil
}

// Send form back to client
func sendForm(w http.ResponseWriter, form interface{}) {
	if buf, err := json.Marshal(form); err != nil {
		log.Printf(logFmt, "sendForm", err)
		return
	} else if _, err := fmt.Fprint(w, string(buf)); err != nil {
		log.Printf(logFmt, "sendForm", err)
	}
}
