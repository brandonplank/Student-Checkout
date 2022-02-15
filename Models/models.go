package models

type Student struct {
	Name    string `csv:"Name"`
	SignOut string `csv:"Signed Out"`
	SignIn  string `csv:"Signed In"`
	Date    string `csv:"Date"`
}
