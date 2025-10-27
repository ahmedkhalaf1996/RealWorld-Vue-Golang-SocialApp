package main

import (
	"Server/realtime"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRealtimeNotifications(app *fiber.App) {
	log.Println("Setting Up Realtime notifications....")

	// get nm
	notificationManager := realtime.GetNotificationManager()

	// n status endpoit
	app.Get("/notifications/status", func(c *fiber.Ctx) error {
		connectedUsers := notificationManager.GetConnectedUsers()
		return c.JSON(fiber.Map{
			"status":          "active",
			"connected_users": len(connectedUsers),
			"users":           connectedUsers,
		})
	})

	// WebSocket route for notficaitons
	app.Get("/notifications/ws/:userId", websocket.New(func(c *websocket.Conn) {
		userID := c.Params("userId")

		log.Printf("New connection Ws for user : %s", userID)

		notificationManager.AddNotificationConnection(userID, c)

		defer func() {
			log.Printf("Closing notifcaton ws conn for user : %s", userID)
			notificationManager.RemoveNotificationConnection(userID)
			c.Close()
		}()

		// keep conn alive and linst to incommeng messges if any
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Printf("Notification WebSocket Error for user %s : %v", userID, err)
				break
			}
		}
	}))
}
