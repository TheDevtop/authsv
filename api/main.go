package api

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheDevtop/authsv/db"
)

const logFmt = "api.%s: %s\n"

// Package entrypoint
func PackageMain() {

	// Declare and parse flags
	flagKeyFile := flag.String("kf", "", "Specify key file")
	flagCertFile := flag.String("cf", "", "Specify certificate file")
	flagDataFile := flag.String("df", "", "Specify database file")
	flag.Parse()

	// Load cached database
	if err := db.LoadCache(*flagDataFile); err != nil {
		log.Printf(logFmt, "PackageMain", err)
		os.Exit(3)
	}

	// Bind authentication functions to DefaultServeMux
	http.HandleFunc("/auth/login", authLogin)
	http.HandleFunc("/auth/list", authList)
	http.HandleFunc("/auth/query", authQuery)

	// Bind administrative functions to DefaultServeMux
	http.HandleFunc("/admin/ping", adminPing)
	http.HandleFunc("/admin/chsecret", adminChangeSecret)
	http.HandleFunc("/admin/user/add", adminUserAdd)
	http.HandleFunc("/admin/user/del", adminUserDel)
	http.HandleFunc("/admin/role/add", adminRoleAdd)
	http.HandleFunc("/admin/role/del", adminRoleDel)

	// Declare and start server
	srv := http.Server{Addr: ":3890", Handler: http.DefaultServeMux}
	go func() {
		if err := srv.ListenAndServeTLS(*flagCertFile, *flagKeyFile); err != nil {
			log.Printf(logFmt, "PackageMain", err)
		}
	}()

	// Wait for signal, and shutdown server
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Printf("api.PackageMain: Caught (%s), stopping server", <-sigCh)
	srv.Shutdown(context.Background())
	if err := db.StoreCache(); err != nil {
		log.Printf(logFmt, "PackageMain", err)
		os.Exit(3)
	}
	os.Exit(0)
}
