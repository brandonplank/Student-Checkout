package models

import "time"

type Main struct {
	AdminName     string   `csv:"Name" json:"adminName" bson:"adminName"`
	AdminPassword string   `csv:"Password" json:"adminPassword" bson:"adminPassword"`
	AdminEmail    string   `csv:"Email" json:"email" bson:"email"`
	Schools       []School `csv:"Schools" json:"schools" bson:"schools"`
}

type School struct {
	Name          string      `csv:"Name" json:"name" bson:"name"`
	Logo          string      `csv:"Logo" json:"logo" bson:"logo"`
	AdminName     string      `csv:"Admin Name" json:"adminName" bson:"adminName"`
	AdminPassword string      `csv:"Admin Password" json:"adminPassword" bson:"adminPassword"`
	AdminEmail    string      `csv:"Email" json:"email" bson:"email"`
	Classrooms    []Classroom `csv:"Classrooms" json:"classrooms" bson:"classrooms"`
}

type Classroom struct {
	Name     string    `csv:"Name" json:"name" bson:"name"`
	Password string    `csv:"Password" json:"password" bson:"password"`
	Email    string    `csv:"Email" json:"email" bson:"email"`
	IsAdmin  bool      `csv:"IsAdmin" json:"isAdmin" bson:"isAdmin"`
	Students []Student `csv:"Students" json:"students" bson:"students"`
}

type PublicStudent struct {
	Name    string `csv:"Name" json:"name" bson:"name"`
	SignOut string `csv:"Signed Out" json:"signedOut" bson:"signedOut"`
	SignIn  string `csv:"Signed In" json:"signedIn" bson:"signedIn"`
	Date    string `csv:"Date" json:"date" bson:"date"`
}
type Student struct {
	Name      string `csv:"Name" json:"name" bson:"name"`
	SignOut   string `csv:"Signed Out" json:"signedOut" bson:"signedOut"`
	SignIn    string `csv:"Signed In" json:"signedIn" bson:"signedIn"`
	Date      string `csv:"Date" json:"date" bson:"date"`
	Classroom string `csv:"Classroom" json:"classroom" bson:"classroom"`
}

type Students []Student
type PublicStudents []PublicStudent

func (p Students) Len() int {
	return len(p)
}

func (p Students) Less(i, j int) bool {
	time1, _ := time.Parse("3:04 pm", p[i].SignOut)
	time2, _ := time.Parse("3:04 pm", p[j].SignOut)
	return time1.Before(time2)
}

func (p Students) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PublicStudents) Len() int {
	return len(p)
}

func (p PublicStudents) Less(i, j int) bool {
	time1, _ := time.Parse("3:04 pm", p[i].SignOut)
	time2, _ := time.Parse("3:04 pm", p[j].SignOut)
	return time1.Before(time2)
}

func (p PublicStudents) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func StudentsToPublicStudents(students Students) PublicStudents {
	var publicStudents PublicStudents
	for _, student := range students {
		publicStudents = append(publicStudents, PublicStudent{Name: student.Name, SignOut: student.SignOut, SignIn: student.SignIn, Date: student.Date})
	}
	return publicStudents
}
