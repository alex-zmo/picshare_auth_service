package main

import (
	"github.com/gmo-personal/picshare_auth_service/database"
	"github.com/gmo-personal/picshare_auth_service/server"
)

func main() {
	database.InitDatabase()
	server.InitServer()
}

