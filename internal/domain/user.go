package domain

import "time"

type User struct {
  Id       int64
  Name     string
  Password string
  Phone    string
  Email    string
  Birthday time.Time
}
