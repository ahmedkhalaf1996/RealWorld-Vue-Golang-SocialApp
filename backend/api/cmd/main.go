package main

import (
	"Server/database"
	"Server/kafka"
	"Server/realtime"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Println("Error loading .env file using environment variables")
	}

	// connect to mongodb database
	database.Connect()

	// Call Redis init connection
	database.InitRedis()
	defer database.CloseRedis()

	// get improtant keys
	nodeID := flag.String("node", "node-1", "Node id for this instance")
	port := flag.String("port", "5000", "port to list on")
	//kafkaAddr := flag.String("kafka", "127.0.0.1:29092", "kafka address")
	kafkaAddr := flag.String("kafka", "127.0.0.1:29092", "kafka address")

	flag.Parse()

	if err := kafka.WaitForKafka(*kafkaAddr, 1*time.Minute); err != nil {
		log.Fatalf("Kafka not ready : %v", err)
	}

	time.Sleep(3 * time.Second)

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))
	// initialize notfication manger with kafka
	if err := realtime.InitNotificationManger(*kafkaAddr, *nodeID); err != nil {
		log.Fatalf("Faild to init notification manager: %v", err)
	}

	// initilize chathub
	chathub := SetupRealtimeChat(app, *kafkaAddr, *nodeID)

	// call heartbeat
	heatbeatMgr, err := kafka.NewHeartbeatManager(*kafkaAddr, *nodeID, chathub)
	if err != nil {
		log.Fatalf("Faild to crete heartbeat manager.. %v", err)
	}

	// Setup API Routes
	SetupAPI(app)
	// Setup realtime Notificatons
	SetupRealtimeNotifications(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Socail app")
	})

	// Serve swager doctionation
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"node":   *nodeID,
		})
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shuting down gracefully..")

		heatbeatMgr.Stop()
		if chathub != nil {
			chathub.Close()
		}

		notifMgr := realtime.GetNotificationManager()
		if notifMgr != nil {
			notifMgr.Close()
		}

		app.ShutdownWithTimeout(10 * time.Second)
	}()
	log.Printf("Server Starting on :%s (node: %s)", *port, *nodeID)
	if err := app.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
