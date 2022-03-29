package main

import (
	"brandonplank.org/checkout/routes"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/mileusna/crontab"
	"log"
	"os"
	"strconv"
	"time"
)

const Port = 8064
const Key = "classof2022"

var context *fiber.Ctx

func Auth(email string, password string) bool {
	if routes.SanitizeString(email) == routes.SanitizeString(routes.MainGlobal.AdminEmail) && password == routes.MainGlobal.AdminPassword {
		return true
	} else {
		for _, school := range routes.MainGlobal.Schools {
			if routes.SanitizeString(email) == routes.SanitizeString(school.AdminEmail) {
				if password == school.AdminPassword {
					return true
				}
			}
			for _, classroom := range school.Classrooms {
				if routes.SanitizeString(classroom.Email) == routes.SanitizeString(email) {
					if password == classroom.Password {
						return true
					}
				}
			}
		}
	}
	return false
}

func setupRoutes(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
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
			ctx.Append("License", "GNU Affero General Public License v3.0")
			context = ctx
			return ctx.Next()
		},
		basicauth.New(basicauth.Config{
			Authorizer:      Auth,
			ContextUsername: "email",
		}),
	)

	app.Static("/", "./Public")

	app.Get("/", routes.Home)
	app.Post("/id/:name", routes.Id)
	app.Post("/search", routes.AdminSearchStudent)
	app.Post("/search/:name", routes.AdminSearchStudent)
	app.Post("/isOut/:name", routes.IsOut)
	app.Get("/GetCSV", routes.GetCSV)
	app.Get("/GetAdminCSV", routes.GetAdminCSV)
	app.Get("/CleanClass", routes.CleanClass)
	app.Get("/CleanClass/:name", routes.CleanClass)
	app.Get("/classroom.csv", routes.CSVFile)
	app.Get("/admin.csv", routes.AdminCSVFile)
	app.Post("/addTeacher", routes.AddTeacher)
	app.Post("/removeTeacher", routes.RemoveTeacher)
	app.Post("/changePassword", routes.ChangePassword)
}

func main() {
	log.Println("[START] Starting student checkout server")

	err := sentry.Init(sentry.ClientOptions{
		Dsn:   "https://f98a7533c0e5433eb0eb89b4f97e5ece@o956450.ingest.sentry.io/6240133",
		Debug: false,
	})
	if err != nil {
		log.Println(err)
	}

	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()
	log.Println("[START] Started Sentry")

	err = godotenv.Load()
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

	ctab.MustAddJob("0 0 * * 1-5", func() { // 12:00 AM every weekday
		routes.CleanStudents()
	})

	engine := html.New("./Resources/Views", ".html")
	router := fiber.New(fiber.Config{DisableStartupMessage: true, Views: engine})
	setupRoutes(router)
	log.Println("[START] Finished setting up routes")

	log.Println("[START] Starting server on port", strconv.Itoa(Port))
	log.Fatalln(router.Listen(":" + strconv.Itoa(Port)))
}
