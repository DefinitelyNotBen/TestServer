package main

import (
	"test/db/database"
	"test/db/server"
)

func main() {
	db := database.NewDB()
	server.Start(db)
}
