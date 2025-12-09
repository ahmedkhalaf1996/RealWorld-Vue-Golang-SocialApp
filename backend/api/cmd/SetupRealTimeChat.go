package main

import (
	"Server/realtime"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupRealtimeChat(app *fiber.App, kafkaAddr, nodeID string) *realtime.ChatHub {

	chatHub, err := realtime.NewChatHub(kafkaAddr, nodeID)
	if err != nil {
		log.Fatalf("Faild to create chat hub %v", err)
	}

	app.Get("/ws/:id", chatHub.HandleClient)

	log.Println("REatime chat Configured successfully with kafka!")
	return chatHub
}
