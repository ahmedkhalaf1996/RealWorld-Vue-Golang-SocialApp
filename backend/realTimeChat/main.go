package main

import (
	"log"
	"realTimeChat/realtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	manager := realtime.NewConnectionManager(realtime.GetUserFriends)
	// register ws route

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		if manager == nil {
			return
		}
		manager.AddConnection(id, c)
		defer func() {
			manager.RemoveConnection(id)
			c.Close()
		}()

		var msg realtime.Message
		for {
			err := c.ReadJSON(&msg)
			if err != nil {
				handleWebSocketError(err, id)
				manager.RemoveConnection(id)
				c.Close()
				break
			}

			log.Printf("Received message form %s to %s : %s", msg.Sender, msg.Recever, msg.Content)
			manager.SendToReceiver(msg)
		}

	}))

	log.Fatal(app.Listen(":8001"))
}

func handleWebSocketError(err error, userID string) {
	log.Printf("WebSocket error for user %s : %v", userID, err)
}
