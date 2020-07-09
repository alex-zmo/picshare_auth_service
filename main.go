package main

import (
	"smart_photos/auth_service/database"
	"smart_photos/auth_service/server"
)

func main() {
	database.InitDatabase()
	server.InitServer()

}

