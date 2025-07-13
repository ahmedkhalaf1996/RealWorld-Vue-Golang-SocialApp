package tests

import (
	"Server/database"
	"Server/routes"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var app *fiber.App

func TestMain(m *testing.M) {
	// Setup
	setup()

	// Run Tests
	code := m.Run()

	// cleanup
	cleanup()

	os.Exit(code)
}

func setup() {
	// Load test environment variables
	if err := godotenv.Load("../.env.test"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
		}
	}

	// Set jws if not prvided
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-jwt-secret-key")
	}

	// Connect to test db
	connectTestDB()

	// setup fiber app
	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// setup routes
	routes.SetupAuthRoutes(app)
	// ..
}

func connectTestDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// use test db
	mongoUri := os.Getenv("TEST_MONGO_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("Faild to connect to test db:", err)
	}
	database.Client = client
	database.DB = client.Database("social_test")
}

func cleanup() {
	if database.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// drop test db
		database.DB.Drop(ctx)
		database.Client.Disconnect(ctx)
	}
}

// helper for cleanup colelctions
func cleanupCollections() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collections := []string{"users"}
	for _, collection := range collections {
		database.DB.Collection(collection).Drop(ctx)
	}
}
