package routes

import (
	"brandonplank.org/checkout/models"
	"encoding/base64"
	"errors"
	csv "github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"os"
	"reflect"
	"time"
)

const csvFileName = "classroom.csv"

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= (valueLen-1)/2; i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func IsStudentOut(name string, students []*models.Student) bool {
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
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	var students = []*models.Student{}

	name := string(nameData)

	studentsFile, err := os.OpenFile(csvFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	defer studentsFile.Close()

	if err := csv.UnmarshalFile(studentsFile, &students); err != nil {
		log.Println("[CSV] No students")
	}

	if err := studentsFile.Truncate(0); err != nil {
		log.Println("[CSV] Unable to truncate start")
	}

	if _, err := studentsFile.Seek(0, 0); err != nil {
		log.Println("[CSV] Unable to seek start")
	}

	if IsStudentOut(name, students) {
		log.Println(name, "has returned")
		var tempStudents []*models.Student
		for _, stu := range students {
			if stu.Name == name {
				if stu.SignIn == "Signed Out" {
					stu.SignIn = time.Now().Format("3:04 pm")
				}
			}
			tempStudents = append(tempStudents, stu)
		}
		students = tempStudents
	} else {
		log.Println(name, "has left")
		students = append(students, &models.Student{Name: name, SignOut: time.Now().Format("3:04 pm"), SignIn: "Signed Out"})
	}

	err = csv.MarshalFile(&students, studentsFile)
	if err != nil {
		panic(err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func GetCSV(ctx *fiber.Ctx) error {
	var students = []*models.Student{}
	studentsFile, err := os.OpenFile(csvFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	defer studentsFile.Close()

	if err := csv.UnmarshalFile(studentsFile, &students); err != nil {
		return ctx.SendString("No students yet")
	}

	ReverseSlice(students)

	content, _ := csv.MarshalBytes(students)

	return ctx.Send(content)
}

func CSVFile(ctx *fiber.Ctx) error {
	return ctx.SendFile(csvFileName, false)
}

func IsOut(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	name := string(nameData)

	var students = []*models.Student{}
	studentsFile, err := os.OpenFile(csvFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	defer studentsFile.Close()

	if err := csv.UnmarshalFile(studentsFile, &students); err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	type out struct {
		IsOut bool   `json:"isOut"`
		Name  string `json:"name"`
	}

	return ctx.JSON(out{IsOut: IsStudentOut(name, students), Name: name})
}

func DoDailyStuff() {

	pass := os.Getenv("PASSWORD")

	e := email.NewEmail()
	e.From = "Brandon Plank <planksprojects@gmail.com>"
	e.To = []string{"susie.hart@rowan.kyschools.us", "brandon@brandonplank.org"}
	e.Subject = "Classroom Sign-Outs"
	e.Text = []byte("This is an automated email")
	e.AttachFile(csvFileName)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "planksprojects@gmail.com", pass, "smtp.gmail.com"))
	if err != nil {
		log.Println(err)
	}

	err = os.Remove(csvFileName)
	if err != nil {
		log.Println("ono")
	}
	studentsFile, err := os.OpenFile(csvFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer studentsFile.Close()
}
