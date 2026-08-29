package server

import (
	"net/http"
	"os"
	"strconv"
)

func StartServer() error {
	port := 7540
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
