package main

import (
	"os"

	"github.com/anargu/miauth"
	"github.com/anargu/miauth/server"
)

func main() {
	// it looks simple:
	// first initialize config params
	miauth.InitConfig()
	// second: initialize DB
	miauth.InitDB()
	// not sure if migration should be inside initDB fn
	miauth.RunMigration()
	defer miauth.CloseDB()

	argsWithoutProg := os.Args[1:]
	// third: run the server. That's it!
	if len(argsWithoutProg) > 0 {
		if argsWithoutProg[0] == "grpc" {
			server.InitGrpcServer()
		} else {
			server.InitServer()
		}
	} else {
		server.InitServer()
	}
}
