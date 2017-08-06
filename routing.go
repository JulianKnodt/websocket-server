package main

import (
  "regexp"
  "fmt"
  "golang.org/x/net/websocket"
)

// General Format is controllerType / Function Name / PayLoad
var routing = map[*regexp.Regexp]func(params string, ws *websocket.Conn){
  regexp.MustCompile("/positions/create/(.*)"): GetPositionController().Create,
  regexp.MustCompile("/positions/show/(.*)"): GetPositionController().Show,
  regexp.MustCompile("/positions/near/(.*)"): GetPositionController().Near,
}

func RouteMessage(msg string, ws *websocket.Conn) {
  for key, matchedFunction := range routing {
    found := key.FindStringSubmatch(msg)
    if found != nil {
      args := found[len(found) - 1]
      matchedFunction(args, ws)
      return
    }
  }
  if err := websocket.Message.Send(ws, "No matching path"); err != nil {
    fmt.Println("No matching path and send failed")
  }
}