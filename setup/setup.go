package setup

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TheDevtop/authsv/db"
)

// Package entrypoint
func PackageMain() {
	flagName := flag.String("n", "admin", "Specify initial administrator")
	flagSecret := flag.String("s", "password", "Specify secret for administrator")
	flagRoles := flag.String("r", db.RoleAdmin, "Specify roles for administrator")
	flagFile := flag.String("f", "", "Specify output file")
	flag.Parse()

	setupDB := make(db.UserDB, 1)
	setupDB[db.UserKey(*flagName)] = struct {
		Secret db.SecretKey
		Roles  db.RoleSet
	}{
		Secret: db.SecretKey(*flagSecret),
		Roles:  db.RoleSet(strings.Fields(*flagRoles)),
	}

	if buf, err := json.Marshal(&setupDB); err != nil {
		fmt.Println(err)
		os.Exit(2)
	} else if err := os.WriteFile(*flagFile, buf, db.FilePerm); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
