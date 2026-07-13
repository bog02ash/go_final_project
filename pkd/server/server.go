package server

import (
	"net/http"
)

func StartServer(port string) {
	if port == "" {
		port = ":7540"
	}
	port = ":" + port
	http.Handle("/", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
