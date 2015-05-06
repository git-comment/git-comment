package main

import (
  "time"
  "github.com/wayn3h0/go-uuid"
)

type Person struct {
  Name string
  Email string
}

type FileRef struct {
  Path string
  Line int
}

type Comment struct {
  Author Person
  CreateTime time.Time
  Content string
  Amender Person
  AmendTime time.Time
  Commit uuid.UUID
  ID uuid.UUID
}

func CurrentPerson() *Person {
  return &Person{}
}

func NewComment(message string, author *Person) *Comment{
  return &Comment{}
}

func main() {

}
