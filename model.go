package main

// For model items
type Item interface {
  Save() (error)
  Delete() (error, Item)
}

type MutableItem interface {
  Update(field, val string) (error, *MutableItem)
}