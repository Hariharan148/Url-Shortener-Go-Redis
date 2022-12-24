package main


import (
	"os"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/Hariharan148/Url-Shortener-Go-Redis/api/routes"
)


func setupRoutes(app *fiber.App){

	app.Get("/:url", routes.ResolveUrl)
	app.Post("/api/v2", routes.ShortenUrl)
}


func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading environment variables ", err)
	}

	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}