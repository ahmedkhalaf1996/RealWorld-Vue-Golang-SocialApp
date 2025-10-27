package main

import (
	"Server/realtime"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRealtimeChat(app *fiber.App) {
	log.Println("Setting Up Realtime chat....")

	// get nm
	manager := realtime.NewConnectionManager(realtime.GetUserFriends)

	// n status endpoit
	app.Get("/realtime/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "realtime-chat",
		})
	})

	// WebSocket route for notficaitons
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		userID := c.Params("id")
		if manager == nil {
			log.Printf("Error: Connection manager is nil for user %s", userID)
		}

		log.Printf("New connection Ws for user : %s", userID)

		manager.AddConnection(userID, c)

		defer func() {
			log.Printf("Closing notifcaton ws conn for user : %s", userID)
			manager.RemoveConnection(userID)
			c.Close()
		}()

		// handle in coming messges
		var msg realtime.Message
		for {
			err := c.ReadJSON(&msg)
			if err != nil {
				handleWebSocketError(err, userID)
				manager.RemoveConnection(userID)
				c.Close()
				break
			}
			log.Printf("Receved message from %s to %s :%s", msg.Sender, msg.Recever, msg.Content)

			manager.SendToReceiver(msg)
		}
	}))

	log.Println("RealTime Chat Configured successfully!")
}

func handleWebSocketError(err error, userID string) {
	log.Printf("WebSocket error for user %s : %v", userID, err)
}
