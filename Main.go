package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

const port = ":3000"

func main() {
	fmt.Printf("Starting server on %s\n", port)
	http.Handle("/api", websocket.Handler(Receive))
	err := http.ListenAndServe(port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
