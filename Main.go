package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
)

const localPort = ":3000"

func main() {
	definedPort := getPort()
	fmt.Printf("Starting server on %s\n", definedPort)
	http.Handle("/api", websocket.Handler(Receive))
	http.HandleFunc("/", placeHolder)
	err := http.ListenAndServe(definedPort, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func placeHolder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello~")
}

func getPort() string {
	osPort := os.Getenv("PORT")
	if osPort != "" {
		return osPort
	}
	return localPort
}
