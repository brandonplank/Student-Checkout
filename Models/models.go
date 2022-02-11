package models

import "time"

type Classmate struct {
	Name    string
	SignOut time.Time
	SignIn  time.Time
}
