package main

import (
  "strconv"
  "encoding/json"
  "fmt"
  "golang.org/x/net/websocket"
  "sync"
  "strings"
)

type positionController struct {}
var controller positionController
var positionOnce sync.Once

func GetPositionController() positionController {
  positionOnce.Do(func() {
    controller = *new(positionController)
  })
  return controller
}

// /api/position/create/"lat:long"
func (pc positionController) Create(params string, ws *websocket.Conn) {
  // todo use regex
  parts := strings.Split(params, ":")
  if len(parts) < 2 {
    websocket.Message.Send(ws, "Invalid params")
    return
  }
  lat := parts[0]
  long := parts[1]
  pos := NewPosition()
  latitudeFloat, err := strconv.ParseFloat(lat, 64)
  pos.Latitude = latitudeFloat
  longitudeFloat, err := strconv.ParseFloat(long, 64)
  pos.Longitude = longitudeFloat
  if err != nil {
    websocket.Message.Send(ws, err.Error())
    return
  }
  pos.Save()
  byteArr, err := json.Marshal(pos)
  result := string(byteArr)
  if err := websocket.Message.Send(ws, result); err != nil {
    fmt.Println(err)
  }
}

// /api/position/show/:uuid
func (pc positionController) Show(uuid string, ws *websocket.Conn) {
  pos, err := FindPosition(uuid)
  if err != nil {
    websocket.Message.Send(ws, err.Error())
  }
  res, err := json.Marshal(pos)
  if err != nil {
    websocket.Message.Send(ws, err.Error())
  }
  if err :=  websocket.Message.Send(ws, string(res)); err != nil {
    fmt.Println(err)
  }
}


// /api/position/near/:uuid
func (pc positionController) Near(uuid string, ws *websocket.Conn) {
  pos, err := FindPosition(uuid)
  if err != nil {
    websocket.Message.Send(ws, err.Error())
    return
  }
  err, positions := pos.FindNearby()

  if err != nil {
    websocket.Message.Send(ws, fmt.Sprintf("Error finding nearby: %s", err.Error()))
    return 
  }
  res, err := json.Marshal(positions)
  if err != nil {
    websocket.Message.Send(ws, fmt.Sprintf("Error marshaling json %S", err.Error()))
    return
  }

  if err := websocket.Message.Send(ws, string(res)); err != nil {
    fmt.Println(err)
  }
}
