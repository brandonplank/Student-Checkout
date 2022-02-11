package routes

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

var people []string

var path = "log.txt"

func CreateLogFile() {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	fmt.Println("File opened at", path)
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func DoesPersonExist(name string) (bool, int) {
	for index, object := range people {
		if name == object {
			return true, index
		}
	}
	return false, 0
}

func Home(ctx *fiber.Ctx) error {
	return ctx.Render("main", fiber.Map{})
}

func Id(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.Next()
	}
	name := string(nameData)
	log.Println(name, nameBase64)

	does, index := DoesPersonExist(name)
	var file, _ = os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()
	if does {
		file.WriteString(fmt.Sprintf("%s\nSigned In at %s\n", name, time.Now().String()))
		people = remove(people, index)
	} else {
		file.WriteString(fmt.Sprintf("%s\nSigned Out at %s\n", name, time.Now().String()))
	}

	return ctx.Next()
}
