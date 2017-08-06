package main

import (
	"github.com/go-redis/redis"
  "sync"
  "errors"
  "reflect"
)

type RedisInstance struct {
  Instance *redis.Client
}

const (
  TYPE_NAME = "type_name"
)

var RedisClient *RedisInstance
var once sync.Once

func GetRedisInstance() *RedisInstance {
  once.Do(func() {
    RedisClient = &RedisInstance{ Instance: redisClient() }
  })
  return RedisClient
}

func (ri RedisInstance) SaveItemType(uuid string, item Item) error {
  return ri.Instance.HMSet(uuid, map[string]interface{}{ 
    TYPE_NAME: reflect.TypeOf(item),
  }).Err()
}

func (ri RedisInstance) GetItemType(uuid string) (typeName string, err error) {
  res, getErr := ri.Instance.HMGet(uuid, TYPE_NAME).Result()
  err = getErr
  if len(res) > 0 {
    resultType, success := res[0].(string)
    typeName = resultType
    if !success {
      err = errors.New("casting item type to string failed")
    }
  } else {
    err = errors.New("No key saved")
  }
  return
}

func redisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

  return client
}
