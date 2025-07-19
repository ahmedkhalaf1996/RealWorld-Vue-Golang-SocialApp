package main

import (
	"Server/database"
	"Server/routes"
	"Server/servergrpc"
	"net"

	pb "Server/protos"
	"log"

	_ "Server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Setup Grpc Server
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("faild to listen : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRealtimeChatServiceServer(grpcServer, &servergrpc.Server{})
	reflection.Register(grpcServer)
	log.Println("gRPC Server Runing on Port 5001")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("faild to server : %v", err)
		}
	}()
	// end of setup gRPC server

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Socail app")
	})

	// setup routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupPostRoutes(app)
	routes.SetupChatRoutes(app)
	routes.SetupNotificationRoutes(app)
	// Serve swager doctionation
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Listen(":5000")
}
