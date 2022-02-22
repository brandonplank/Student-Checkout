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
	"log"
	"os"
	"strconv"
	"strings"
)

const Port = 8064
const Key = "classof2022"

var context *fiber.Ctx

func Auth(name string, password string) bool {
	if name == routes.MainGlobal.AdminName && password == routes.MainGlobal.AdminPassword {
		return true
	} else {
		for _, school := range routes.MainGlobal.Schools {
			for _, classroom := range school.Classrooms {
				if strings.ToLower(classroom.Name) == strings.ToLower(name) {
					if password == classroom.Password {
						return true
					}
					return false
				}
			}
		}
	}

	//if context == nil {
	//	log.Println("That's not supposed to happen")
	//	return false
	//}
	//cookie := context.Cookies("token")
	//
	//if len(cookie) > 5 {
	//	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
	//		return []byte(Key), nil
	//	})
	//	if err != nil {
	//		context.Status(fiber.StatusUnauthorized)
	//	}
	//	claims := token.Claims.(*jwt.StandardClaims)
	//
	//	err = claims.Valid()
	//	if err != nil {
	//		// destroy token
	//		context.Cookie(&fiber.Cookie{
	//			Name:     "token",
	//			Value:    "",
	//			Expires:  time.Now().Add(-(time.Hour * 2)),
	//			HTTPOnly: true,
	//		})
	//		context.Status(fiber.StatusUnauthorized)
	//	}
	//
	//	if claims.Issuer == name {
	//		return true
	//	}
	//	return false
	//}
	//
	//if name == routes.MainGlobal.AdminName && password == routes.MainGlobal.AdminPassword {
	//	return true
	//} else {
	//	for _, school := range routes.MainGlobal.Schools {
	//		for _, classroom := range school.Classrooms {
	//			if strings.ToLower(classroom.Name) == strings.ToLower(name) {
	//				err := bcrypt.CompareHashAndPassword(classroom.Password, []byte(password))
	//				if err != nil {
	//					return false
	//				}
	//				claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
	//					Issuer:    name,
	//					ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	//				})
	//
	//				token, err := claims.SignedString([]byte(Key))
	//				if err != nil {
	//					log.Println(err)
	//				}
	//
	//				context.Cookie(&fiber.Cookie{
	//					Name:     "token",
	//					Value:    token,
	//					Expires:  time.Now().Add(24 * time.Hour),
	//					HTTPOnly: true,
	//				})
	//				return true
	//			}
	//		}
	//	}
	//}
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
			ctx.Append("License", "BSD 3-Clause License")
			context = ctx
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
	app.Get("/CleanClass", routes.CleanClass)
	app.Get("/CleanClass/:name", routes.CleanClass)
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
