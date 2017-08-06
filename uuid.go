package main

import (
  "github.com/nu7hatch/gouuid"
)

// uuid is [16]byte

func GenerateUuid() (string) {
  uuid, err := uuid.NewV4()
  if err != nil {
    panic(err)
  }
  return uuid.String()
}
