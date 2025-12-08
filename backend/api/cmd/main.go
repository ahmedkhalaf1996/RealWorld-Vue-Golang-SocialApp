package main

import (
	"Server/database"
	"flag"

	"log"

	_ "Server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Fiber Golang Mongo Grpc Websocet etc..
// @version 1.0
// @description This is Swagger docs for rest api golang fiber
// @host localhost:5000
// @BasePath /
// @schemes http
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token

func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := flag.String("port", "5000", "Port to listen on")
	flag.Parse()
	// connect to mongodb database
	database.Connect()

	// Call Redis init connection
	database.InitRedis()
	defer database.CloseRedis()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Setup API Routes
	SetupAPI(app)

	// setup realtime chat
	SetupRealtimeChat(app)

	// Setup realtime Notificatons
	SetupRealtimeNotifications(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Socail app")
	})

	// Serve swager doctionation
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Printf("Server Starting on :%s", *port)
	if err := app.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
