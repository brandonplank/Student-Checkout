package main

import (
	"brandonplank.org/checkout/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/mileusna/crontab"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
)

const Port = 8064

func Auth(name string, password string) bool {
	if name == routes.MainGlobal.AdminName && password == routes.MainGlobal.AdminPassword {
		return true
	} else {
		for _, school := range routes.MainGlobal.Schools {
			for _, classroom := range school.Classrooms {
				if classroom.Name == name {
					err := bcrypt.CompareHashAndPassword([]byte(classroom.Password), []byte(password))
					if err != nil {
						return false
					}
					return true
				}
			}
		}
	}
	return false
}

func setupRoutes(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowHeaders: "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
		}),
		logger.New(logger.Config{
			Format:     "${time} [${method}]->${status} Latency->${latency} - ${path} | ${error}\n",
			TimeFormat: "2006/01/02 15:04:05",
		}),
		cors.New(cors.Config{
			AllowCredentials: true,
		}),
		func(ctx *fiber.Ctx) error {
			ctx.Append("Access-Control-Allow-Origin", "*")
			ctx.Append("Developer", "Brandon Plank")
			ctx.Append("License", "BSD 3-Clause License")
			return ctx.Next()
		},
		basicauth.New(basicauth.Config{
			Authorizer:      Auth,
			ContextUsername: "name",
		}),
	)

	//serve := app.Group("/assets")
	app.Static("/", "./Public")

	app.Get("/", routes.Home)
	app.Post("/id/:name", routes.Id)
	app.Post("/isOut/:name", routes.IsOut)
	app.Get("/GetCSV", routes.GetCSV)
	app.Get("/CleanJSON", routes.CleanJSON)
	app.Get("/classroom.csv", routes.CSVFile)
}

func main() {
	log.Println("[START] Starting student checkout server")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("[ERROR] Error loading .env file")
	}

	database, err := os.OpenFile(routes.DatabaseFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	routes.ReadJSONToStruct()

	ctab := crontab.New()

	ctab.MustAddJob("5 15 * * 1-5", func() { // 03:05 PM every weekday
		routes.DailyRoutine()
	})

	engine := html.New("./Resources/Views", ".html")
	router := fiber.New(fiber.Config{DisableStartupMessage: true, Views: engine})
	setupRoutes(router)
	log.Println("[START] Finished setting up routes")

	log.Println("[START] Starting server on port", strconv.Itoa(Port))
	log.Fatalln(router.Listen(":" + strconv.Itoa(Port)))
}
