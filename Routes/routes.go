package routes

import (
	"brandonplank.org/checkout/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	csv "github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"reflect"
	"sort"
	"sync"
	"time"
)

var MainGlobal = new(models.Main)

const DatabaseFile = "Storage/database.json"
const csvFileName = "classroom.csv"

var mutex sync.Mutex

func WriteJSONToFile() {
	database, err := os.OpenFile(DatabaseFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	data, err := json.MarshalIndent(MainGlobal, "", "\t")

	err = ioutil.WriteFile(DatabaseFile, data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadJSONToStruct() {
	content, _ := ioutil.ReadFile(DatabaseFile)
	if len(content) <= 1 {
		mainModel, _ := json.Marshal(models.Main{})
		err := ioutil.WriteFile(DatabaseFile, mainModel, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := json.Unmarshal(content, &MainGlobal)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	if valueLen < 1 {
		return
	}
	for i := 0; i <= (valueLen-1)/2; i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func IsStudentOut(name string, students []models.Student) bool {
	if students == nil {
		return false
	}
	for _, stu := range students {
		if stu.Name == name {
			if stu.SignIn == "Signed Out" {
				return true
			}
		}
	}
	return false
}

func Home(ctx *fiber.Ctx) error {
	return ctx.Render("main", fiber.Map{"year": time.Now().Format("2006")})
}

func Id(ctx *fiber.Ctx) error {
	name := ctx.Locals("name")

	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	studentName := string(nameData)

	for schoolIndex, school := range MainGlobal.Schools {
		for classroomIndex, classroom := range school.Classrooms {
			if classroom.Name == name {
				if IsStudentOut(studentName, classroom.Students) {
					log.Println(studentName, "has returned")
					var tempStudents []models.Student
					for _, stu := range classroom.Students {
						if stu.Name == studentName {
							if stu.SignIn == "Signed Out" {
								stu.SignIn = time.Now().Format("3:04 pm")
							}
						}
						tempStudents = append(tempStudents, stu)
					}
					mutex.Lock()
					MainGlobal.Schools[schoolIndex].Classrooms[classroomIndex].Students = tempStudents
					mutex.Unlock()
				} else {
					log.Println(studentName, "has left")
					mutex.Lock()
					MainGlobal.Schools[schoolIndex].Classrooms[classroomIndex].Students = append(classroom.Students, models.Student{Name: studentName, SignOut: time.Now().Format("3:04 pm"), SignIn: "Signed Out", Date: time.Now().Format("01/02/2006")})
					mutex.Unlock()
				}
				WriteJSONToFile()
				return ctx.SendStatus(fiber.StatusOK)
			}
		}
	}
	return ctx.SendStatus(fiber.StatusBadRequest)
}

func GetCSV(ctx *fiber.Ctx) error {
	name := ctx.Locals("name")
	for _, school := range MainGlobal.Schools {
		if len(school.Classrooms) > 0 {
			for _, classroom := range school.Classrooms {
				if classroom.Name == name {
					if len(classroom.Students) < 1 {
						return ctx.SendString("No students yet")
					}
					sort.Slice(classroom.Students, func(i, j int) bool {
						time1, _ := time.Parse("01/02/2006", classroom.Students[i].SignOut)
						time2, _ := time.Parse("01/02/2006", classroom.Students[j].SignOut)
						return time1.Before(time2)
					})
					ReverseSlice(classroom.Students)
					content, _ := csv.MarshalBytes(classroom.Students)
					return ctx.Send(content)
				}
			}
		}
	}
	return ctx.SendStatus(fiber.StatusInternalServerError)
}

func CSVFile(ctx *fiber.Ctx) error {
	name := ctx.Locals("name")
	for _, school := range MainGlobal.Schools {
		for _, classroom := range school.Classrooms {
			if classroom.Name == name {
				sort.Slice(classroom.Students, func(i, j int) bool {
					time1, _ := time.Parse("01/02/2006", classroom.Students[i].SignOut)
					time2, _ := time.Parse("01/02/2006", classroom.Students[j].SignOut)
					return time1.Before(time2)
				})
				students, err := csv.MarshalBytes(classroom.Students)
				if err != nil {
					return ctx.SendStatus(fiber.StatusBadRequest)
				}
				ctx.Append("Content-Disposition", "attachment; filename=\"classroom.csv\"")
				ctx.Append("Content-Type", "text/csv")
				return ctx.Send(students)
			}
		}
	}
	return ctx.SendStatus(fiber.StatusBadRequest)
}

func IsOut(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	studentName := string(nameData)

	name := ctx.Locals("name")
	for _, school := range MainGlobal.Schools {
		for _, classroom := range school.Classrooms {
			if classroom.Name == name {
				type out struct {
					IsOut bool   `json:"isOut"`
					Name  string `json:"name"`
				}
				return ctx.JSON(out{IsOut: IsStudentOut(studentName, classroom.Students), Name: studentName})
			}
		}
	}
	return ctx.SendStatus(fiber.StatusBadRequest)
}

func CleanJSON(ctx *fiber.Ctx) error {
	err := os.Remove(DatabaseFile)
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func DailyRoutine() {

	pass := os.Getenv("PASSWORD")

	studentsFile, err := os.OpenFile(DatabaseFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer studentsFile.Close()

	for _, school := range MainGlobal.Schools {
		for _, class := range school.Classrooms {
			students := class.Students

			csvClass, err := csv.MarshalBytes(students)
			if err != nil {
				log.Println(err)
			}
			if len(csvClass) < 5 {
				continue
			}
			csvReader := bytes.NewReader(csvClass)
			e := email.NewEmail()
			e.From = "Brandon Plank <planksprojects@gmail.com>"
			e.To = []string{class.Email}
			e.Subject = "Classroom Sign-Outs"
			e.Text = []byte("This is an automated email to " + class.Name)
			e.Attach(csvReader, fmt.Sprintf("%s.csv", class.Name), "text/csv")
			err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "planksprojects@gmail.com", pass, "smtp.gmail.com"))
			if err != nil {
				log.Println(err)
			}
		}
	}

	err = os.Remove(csvFileName)
	if err != nil {
		log.Println("ono")
	}
}
