package main

import (
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

	// third: run the server. That's it!
	server.InitServer()
}
