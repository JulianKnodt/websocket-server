package main

import (
  "errors"
)

type Person struct {
  Uid string `json: uid`
  Name string `json: name`
  Emote string `json: emote`
}

const (
  PERSON = "person"
  NAME = "name"
  EMOTE = "emote"
)

func (person Person) TypeName() string {
  return PERSON
}

func (person Person) toMap() (result map[string]interface{}) {
  result[NAME] = person.Name
  result[EMOTE] = person.Emote
  return
}

func (person *Person) Save() (*Person, error) {
  redisInstance := GetRedisInstance()
  if person.Uid == "" {
    uuid := GenerateUuid()
    person.Uid = uuid

    if _, err := redisInstance.Instance.HMSet(uuid, person.toMap()).Result(); err != nil {
      return nil, err
    }
  } else {
    if _, err := redisInstance.Instance.HMSet(person.Uid, person.toMap()).Result(); err != nil {
      return nil, err
    }
  }
  return person, nil
}

func (person *Person) Delete() (*Person, error) {
  if person.Uid != "" {
    redisInstance := GetRedisInstance()
    redisInstance.Instance.Del(person.Uid)
    person.Uid = ""

    return person, nil
  } else {
    return person, errors.New("No Uid set on current person")
  }
}

func (person *Person) Find(uid string) (*Person, error) {
  redisInstance := GetRedisInstance()
  result, err := redisInstance.Instance.HGetAll(person.Uid).Result()
  if err != nil {
    return person, nil
  }
  if person == nil {
     person = &Person{}
  }
  person.Uid = uid
  person.Name = result["name"]
  person.Emote = result["emote"]
  return person, err
}
