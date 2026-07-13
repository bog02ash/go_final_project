package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	"my_project/pkd/api"
	"my_project/pkd/db"
	"my_project/pkd/server"
)

func main() {
	_ = godotenv.Load()
	dbFile := os.Getenv("TODO_DBFILE")
	port := os.Getenv("TODO_PORT")
	password := os.Getenv("TODO_PASSWORD")
	api.InitPassword(password)
	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}
	defer db.CloseDB()
	api.Init()
	server.StartServer(port)
}
