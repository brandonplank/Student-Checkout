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
	"strings"
	"sync"
	"time"
)

var MainGlobal = new(models.Main)

const DatabaseFile = "Storage/database.json"

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

func IsAdmin(name string) bool {
	for _, school := range MainGlobal.Schools {
		if strings.ToLower(name) == strings.ToLower(school.AdminName) {
			return true
		}
	}
	return false
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
	name := ctx.Locals("name")
	logoURL := "assets/img/viking_logo.png"
	for _, school := range MainGlobal.Schools {
		for _, classroom := range school.Classrooms {
			if classroom.Name == name {
				if len(school.Logo) > 0 {
					logoURL = school.Logo
					break
				}
			}
		}
	}

	if IsAdmin(name.(string)) {
		return ctx.Render("admin", fiber.Map{
			"year": time.Now().Format("2006"),
			"logo": logoURL,
		})
	}

	return ctx.Render("main", fiber.Map{
		"year": time.Now().Format("2006"),
		"logo": logoURL,
	})
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
					MainGlobal.Schools[schoolIndex].Classrooms[classroomIndex].Students = append(classroom.Students, models.Student{Name: studentName, SignOut: time.Now().Format("3:04 pm"), SignIn: "Signed Out", Date: time.Now().Format("01/02/2006"), Classroom: classroom.Name})
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
					var students models.PublicStudents
					students = models.StudentsToPublicStudents(classroom.Students)
					sort.Sort(students)
					ReverseSlice(students)
					content, _ := csv.MarshalBytes(students)
					return ctx.Send(content)
				}
			}
		}
	}
	return ctx.SendStatus(fiber.StatusInternalServerError)
}

func GetAdminCSV(ctx *fiber.Ctx) error {
	for _, school := range MainGlobal.Schools {
		if len(school.Classrooms) > 0 {
			var allStudents models.Students
			for _, classroom := range school.Classrooms {
				if len(classroom.Students) < 1 {
					continue
				}
				for _, student := range classroom.Students {
					allStudents = append(allStudents, student)
				}
			}
			sort.Sort(allStudents)
			ReverseSlice(allStudents)
			content, _ := csv.MarshalBytes(allStudents)
			return ctx.Send(content)
		}
	}
	return ctx.SendStatus(fiber.StatusInternalServerError)
}

func AdminSearchStudent(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	studentName := string(nameData)

	for _, school := range MainGlobal.Schools {
		if len(school.Classrooms) > 0 {
			var allStudents models.Students
			for _, classroom := range school.Classrooms {
				if len(classroom.Students) < 1 {
					continue
				}
				for _, student := range classroom.Students {
					if strings.Contains(strings.ToLower(student.Name), strings.ToLower(studentName)) {
						allStudents = append(allStudents, student)
					}
				}
			}
			sort.Sort(allStudents)
			ReverseSlice(allStudents)
			content, _ := csv.MarshalBytes(allStudents)
			return ctx.Send(content)
		}
	}
	return ctx.SendStatus(fiber.StatusInternalServerError)
}

func CSVFile(ctx *fiber.Ctx) error {
	name := ctx.Locals("name")
	for _, school := range MainGlobal.Schools {
		for _, classroom := range school.Classrooms {
			if classroom.Name == name {
				var students models.PublicStudents
				students = models.StudentsToPublicStudents(classroom.Students)
				sort.Sort(students)
				studentsBytes, err := csv.MarshalBytes(students)
				if err != nil {
					return ctx.SendStatus(fiber.StatusBadRequest)
				}
				ctx.Append("Content-Disposition", "attachment; filename=\"classroom.csv\"")
				ctx.Append("Content-Type", "text/csv")
				return ctx.Send(studentsBytes)
			}
		}
	}
	return ctx.SendStatus(fiber.StatusBadRequest)
}

func AdminCSVFile(ctx *fiber.Ctx) error {
	for _, school := range MainGlobal.Schools {
		if len(school.Classrooms) > 0 {
			var allStudents models.Students
			for _, classroom := range school.Classrooms {
				if len(classroom.Students) < 1 {
					continue
				}
				for _, student := range classroom.Students {
					allStudents = append(allStudents, student)
				}
			}
			sort.Sort(allStudents)
			ReverseSlice(allStudents)
			content, _ := csv.MarshalBytes(allStudents)
			ctx.Append("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.csv\"", school.Name))
			ctx.Append("Content-Type", "text/csv")
			return ctx.Send(content)
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

func CleanClass(ctx *fiber.Ctx) error {
	var classroomName string
	name := ctx.Locals("name")
	nameBase64 := ctx.Params("name")
	if len(nameBase64) > 0 {
		nameData, err := base64.URLEncoding.DecodeString(nameBase64)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		classroomName = string(nameData)
	}

	if len(classroomName) < 1 {
		classroomName = name.(string)
	}

	for schoolsIndex, school := range MainGlobal.Schools {
		for classroomsIndex, classroom := range school.Classrooms {
			if classroom.Name == classroomName {
				mutex.Lock()
				MainGlobal.Schools[schoolsIndex].Classrooms[classroomsIndex].Students = models.Students{}
				mutex.Unlock()
				WriteJSONToFile()
				return ctx.SendStatus(fiber.StatusOK)
			}
		}
	}
	return ctx.SendStatus(fiber.StatusNotFound)
}

func DoesSchoolHaveStudents(classes []models.Classroom) bool {
	for _, class := range classes {
		if len(class.Students) > 0 {
			return true
		}
	}
	return false
}

func CleanStudents() {
	for schoolsIndex, school := range MainGlobal.Schools {
		for classroomsIndex := range school.Classrooms {
			MainGlobal.Schools[schoolsIndex].Classrooms[classroomsIndex].Students = models.Students{}
		}
	}
}

func AddTeacher(ctx *fiber.Ctx) error {
	var payload map[string]interface{}
	err := ctx.BodyParser(&payload)
	if err != nil {
		return err
	}
	//TeacherName := payload["name"]
	return nil
}

func DailyRoutine() {
	pass := os.Getenv("PASSWORD")

	studentsFile, _ := os.OpenFile(DatabaseFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer studentsFile.Close()

	for _, school := range MainGlobal.Schools {
		if len(school.AdminEmail) < 1 || len(school.AdminName) < 1 || len(school.AdminPassword) < 1 {
			continue
		}
		if DoesSchoolHaveStudents(school.Classrooms) {
			var allStudents models.Students
			for _, classroom := range school.Classrooms {
				if len(classroom.Students) < 1 {
					continue
				}
				for _, student := range classroom.Students {
					allStudents = append(allStudents, student)
				}
			}
			sort.Sort(allStudents)
			ReverseSlice(allStudents)
			content, _ := csv.MarshalBytes(allStudents)
			csvReader := bytes.NewReader(content)

			schoolEmail := email.NewEmail()
			schoolEmail.From = "Classroom Attendance <planksprojects@gmail.com>"
			schoolEmail.Subject = "Classroom Sign-Outs"
			schoolEmail.To = []string{school.AdminEmail}
			schoolEmail.Text = []byte("This is an automated email to " + school.Name)
			schoolEmail.Attach(csvReader, fmt.Sprintf("%s.csv", school.Name), "text/csv")
			err := schoolEmail.Send("smtp.gmail.com:587", smtp.PlainAuth("", "planksprojects@gmail.com", pass, "smtp.gmail.com"))
			if err != nil {
				log.Println(err)
			}
		}
	}

	for _, school := range MainGlobal.Schools {
		for _, class := range school.Classrooms {
			students := class.Students
			if len(students) < 1 {
				continue
			}
			csvClass, err := csv.MarshalBytes(students)
			if err != nil {
				log.Println(err)
			}
			if len(csvClass) < 5 {
				continue
			}
			csvReader := bytes.NewReader(csvClass)
			classroomEmail := email.NewEmail()
			classroomEmail.From = "Classroom Attendance <planksprojects@gmail.com>"
			classroomEmail.Subject = "Classroom Sign-Outs"
			classroomEmail.To = []string{class.Email}
			classroomEmail.Text = []byte("This is an automated email to " + class.Name)
			classroomEmail.Attach(csvReader, fmt.Sprintf("%s.csv", class.Name), "text/csv")
			err = classroomEmail.Send("smtp.gmail.com:587", smtp.PlainAuth("", "planksprojects@gmail.com", pass, "smtp.gmail.com"))
			if err != nil {
				log.Println(err)
			}
		}
	}
	WriteJSONToFile()
}
