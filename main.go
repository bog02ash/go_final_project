package main

import (
	"my_project/pkd/api"
	"my_project/pkd/db"
	"my_project/pkd/server"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	_ = godotenv.Load()
	dbFile := os.Getenv("TODO_DBFILE")
	port := os.Getenv("TODO_PORT")
	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}
	defer db.CloseDB()
	api.Init()
	server.StartServer(port)
}
