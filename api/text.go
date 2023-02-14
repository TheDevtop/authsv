package api

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheDevtop/authsv/db"
)

// Ping server
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

// Check if user and secret matches
func authLogin(w http.ResponseWriter, r *http.Request) {
	const funcName = "authLogin"
	var form formLogin
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if rs, err := db.Login(form.User, form.Secret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else {
		sendForm(w, formReply{Error: "", Data: rs})
	}
}

// Check if user matches role
func authQuery(w http.ResponseWriter, r *http.Request) {
	const funcName = "authQuery"
	var form formQuery
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.Query(form.User, form.Role); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else {
		sendForm(w, formReply{Error: "", Data: form.Role})
	}
}

// Dump entire cachedDB
func authDump(w http.ResponseWriter, r *http.Request) {
	const funcName = "authDump"
	var form formLogin
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AdminLogin(form.User, form.Secret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	}
	sendForm(w, formReply{Error: "", Data: db.Dump()})
}

// Change user secret
func authChangeSecret(w http.ResponseWriter, r *http.Request) {
	const funcName = "authChangeSecret"
	var form formModifyUser
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.ChangeSecret(form.User, form.Secret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	}
	sendForm(w, formReply{Error: "", Data: nil})
}

// Add new user
func authAddUser(w http.ResponseWriter, r *http.Request) {
	const funcName = "authAddUser"
	var form formModifyUser
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AddUser(form.User, form.Secret, form.Roles); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	}
	sendForm(w, formReply{Error: "", Data: nil})
}

// Delete user
func authDeleteUser(w http.ResponseWriter, r *http.Request) {
	const funcName = "authDeleteUser"
	var form formModifyUser
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.DelUser(form.User); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	}
	sendForm(w, formReply{Error: "", Data: nil})
}

// Add role to user
func authAddRole(w http.ResponseWriter, r *http.Request) {
	const funcName = "authAddRole"
	var form formModifyRole
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AddRole(form.User, form.Role); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	}
	sendForm(w, formReply{Error: "", Data: nil})
}

// Remove role from user
func authDeleteRole(w http.ResponseWriter, r *http.Request) {
	const funcName = "authDeleteRole"
	var form formModifyRole
	if err := receiveForm(r, &form); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.AdminLogin(form.AdminUser, form.AdminSecret); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	} else if err := db.DelRole(form.User, form.Role); err != nil {
		log.Printf(logFormat, funcName, err)
		sendForm(w, formReply{Error: err.Error(), Data: nil})
		return
	}
	sendForm(w, formReply{Error: "", Data: nil})
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
		log.Printf(logFormat, "sendForm", err)
		return
	} else if _, err := fmt.Fprint(w, string(buf)); err != nil {
		log.Printf(logFormat, "sendForm", err)
	}
}

// Package entrypoint
func PackageMain() {

	// Declare and parse flags
	flagKeyFile := flag.String("kf", "", "Specify key file")
	flagCertFile := flag.String("cf", "", "Specify certificate file")
	flagDataFile := flag.String("df", "", "Specify database file")
	flag.Parse()

	// Load cached database
	if err := db.LoadCache(*flagDataFile); err != nil {
		log.Printf(logFormat, "PackageMain", err)
		os.Exit(3)
	}

	// Regular functions
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/auth/login", authLogin)
	http.HandleFunc("/auth/query", authQuery)

	// Administrative functions
	http.HandleFunc("/auth/dump", authDump)
	http.HandleFunc("/auth/changeSecret", authChangeSecret)
	http.HandleFunc("/auth/addUser", authAddUser)
	http.HandleFunc("/auth/deleteUser", authDeleteUser)
	http.HandleFunc("/auth/addRole", authAddRole)
	http.HandleFunc("/auth/deleteRole", authDeleteRole)

	// Declare and start server
	srv := http.Server{Addr: ":3890", Handler: http.DefaultServeMux}
	go func() {
		if err := srv.ListenAndServeTLS(*flagCertFile, *flagKeyFile); err != nil {
			log.Printf(logFormat, "PackageMain", err)
		}
	}()

	// Wait for signal, and shutdown server
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Printf("api.PackageMain: Caught (%s), stopping server", <-sigCh)
	srv.Shutdown(context.Background())
	if err := db.StoreCache(); err != nil {
		log.Printf(logFormat, "PackageMain", err)
		os.Exit(3)
	}
	os.Exit(0)
}
