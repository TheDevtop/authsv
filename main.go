package main

import (
	"os"

	"github.com/TheDevtop/authsv/api"
	"github.com/TheDevtop/authsv/setup"
)

// Program entrypoint
func main() {
	const (
		usage   = "usage: authsv server|setup [options...]"
		exitErr = 1
	)
	if len(os.Args) < 2 {
		println(usage)
		os.Exit(exitErr)
	} else if os.Args[1] == "server" {
		os.Args = os.Args[1:]
		api.PackageMain()
	} else if os.Args[1] == "setup" {
		os.Args = os.Args[1:]
		setup.PackageMain()
	} else {
		println(usage)
		os.Exit(exitErr)
	}
}
