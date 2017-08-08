package main

import (
  "golang.org/x/net/websocket"
  "fmt"
)

const (
  ping = "ping"
  pong = "pong"
  default_error = "cannot handle message"
  valid_client = "TEMP"
)

func Receive(ws *websocket.Conn) {
    var err error

    valid := validate(ws)

    if !valid {
      websocket.Message.Send(ws, "Invalid client origin")
      return
    }

    var msg string

    if err = websocket.Message.Receive(ws, &msg); err != nil {
        fmt.Printf("Error receiving message: %s", err.Error())
        return
    }

    go handle(msg, ws)
}

func handle(msg string, ws *websocket.Conn) {
  fmt.Println(msg)
  if msg == ping {
    if err := websocket.Message.Send(ws, pong); err != nil {
      panic(err)
    }
  } else {
    RouteMessage(msg, ws)
  }
}

func validate(ws *websocket.Conn) (valid bool) {
  connectionConfig := ws.Config()
  return connectionConfig.Origin.String() == valid_client || true
}
