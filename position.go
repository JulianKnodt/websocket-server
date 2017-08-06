package main

import (
  "github.com/go-redis/redis"
  "errors"
  "fmt"
  "strconv"
)

type position struct {
  Latitude float64
  Longitude float64
  Region string
  Uuid string
}

const (
  UUID = "uuid"
  REGION = "reg"
)

const GLOBAL_TEMP = "gl_temp"

func NewPosition() *position {
  position := new(position)
  position.Uuid = GenerateUuid()
  position.Region = GLOBAL_TEMP
  return position
}

func toString(fl float64) string {
  return strconv.FormatFloat(fl, 'f', -1, 64)
}

func FindPosition(uuid string) (result *position, err error) {
  redisInstance := GetRedisInstance().Instance
  res, err := redisInstance.GeoPos(GLOBAL_TEMP, uuid).Result();
  if err != nil {
    return nil, err
  }
  if len(res) > 0 {
    geoLocation := res[0]
    result := new(position)
    result.Latitude = geoLocation.Latitude
    result.Longitude = geoLocation.Longitude
    result.Region = GLOBAL_TEMP
    result.Uuid = uuid
    return result, nil
  } else {
    err = errors.New(fmt.Sprint("No position found for uuid %s", uuid))
    return nil, err
  }
}

func Update(uuid string, updated position) (err error, result *position) {
  old, err := FindPosition(uuid)
  if err != nil {
    return
  }
  redisInstance := GetRedisInstance().Instance
  err = redisInstance.SMove(toString(old.Latitude), toString(updated.Latitude), uuid).Err()
  err = redisInstance.SMove(toString(old.Longitude), toString(updated.Longitude), uuid).Err()
  err = redisInstance.GeoAdd(updated.Region, updated.toGeoLocation()).Err()

  return err, &updated
}

func (position position) toGeoLocation() *redis.GeoLocation {
  location := redis.GeoLocation{
    Name: position.Uuid,
    Latitude: position.Latitude,
    Longitude: position.Longitude,
  }
  return &location
}

func (position position) toMap() map[string]interface{} {
  out := map[string]interface{}{}
  out[REGION] = position.Region
  out[UUID] = position.Uuid
  return out
}

func (position position) Save() (err error) {
  redisInstance := GetRedisInstance().Instance

  geoLocation := position.toGeoLocation()

  err = redisInstance.GeoAdd(position.Region, geoLocation).Err()
  err = redisInstance.HMSet(position.Uuid, position.toMap()).Err()
  err = redisInstance.SAdd(fmt.Sprintf("lat%s", toString(position.Latitude)), position.Uuid).Err()
  err = redisInstance.SAdd(fmt.Sprintf("long%s", toString(position.Longitude)), position.Uuid).Err()
  return
}

func (receiver position) Delete() (err error, old *position) {
  redisInstance := GetRedisInstance().Instance
  err = redisInstance.ZRem(receiver.Region, receiver.Uuid).Err()
  err = redisInstance.Del(receiver.Uuid).Err()
  err = redisInstance.SRem(toString(receiver.Latitude), receiver.Uuid).Err()
  err = redisInstance.SRem(toString(receiver.Longitude), receiver.Uuid).Err()
  return err, &receiver
}

func (receiver position) FindNearby() (err error, positions []*position) {
  redisInstance := GetRedisInstance().Instance
  geoLocationQuery := &redis.GeoRadiusQuery{
    Radius: 5,
    Count: 50,
    WithCoord: true,
  }
  res, err := redisInstance.GeoRadiusByMember(GLOBAL_TEMP, receiver.Uuid, geoLocationQuery).Result()
  if err != nil {
    return
  }
  for _, geoLocation := range res {
    returnedPosition := position{
      Longitude: geoLocation.Longitude,
      Latitude: geoLocation.Latitude,
      Region: GLOBAL_TEMP,
    }

    latString := fmt.Sprintf("lat%s", toString(returnedPosition.Latitude))
    longString := fmt.Sprintf("long%s", toString(returnedPosition.Longitude))
    uuids, intersectError := redisInstance.SInter(latString, longString).Result()

    if intersectError != nil {
      err = intersectError
      return
    }

    fmt.Println(uuids)

    if len(uuids) > 0 {
      returnedPosition.Uuid = uuids[0]
    }
    positions = append(positions, &returnedPosition)
  }
  return
}
