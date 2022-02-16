package models

type Main struct {
	AdminName     string   `csv:"Name" json:"adminName"`
	AdminPassword string   `csv:"Password" json:"adminPassword"`
	AdminEmail    string   `csv:"Email" json:"email"`
	Schools       []School `csv:"Schools" json:"schools"`
}

type School struct {
	AdminName     string      `csv:"Name" json:"adminName"`
	AdminPassword string      `csv:"Password" json:"adminPassword"`
	AdminEmail    string      `csv:"Email" json:"email"`
	Classrooms    []Classroom `csv:"Classrooms" json:"classrooms"`
}

type Classroom struct {
	Name     string    `csv:"Name" json:"name"`
	Password string    `csv:"Password" json:"password"`
	Email    string    `csv:"Email" json:"email"`
	Students []Student `csv:"Students" json:"students"`
}

type Student struct {
	Name    string `csv:"Name" json:"name"`
	SignOut string `csv:"Signed Out" json:"signedOut"`
	SignIn  string `csv:"Signed In" json:"signedIn"`
	Date    string `csv:"Date" json:"date"`
}
