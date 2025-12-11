package realtime

import (
	"Server/kafka"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (h *ChatHub) HandleClient(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		userID := c.Params("id")
		if userID == "" {
			return fiber.ErrBadRequest
		}

		return websocket.New(func(conn *websocket.Conn) {
			conn.SetReadLimit(8192)

			client := h.RegisterClient(userID, conn)

			ctx, cancel := context.WithCancel(context.Background())

			cleanupDone := false

			cleanup := func() {
				if !cleanupDone {
					cleanupDone = true
					h.UnregisterClient(userID)
					conn.Close()
					cancel()
				}
			}
			defer cleanup()
			go h.handlePingPong(ctx, client, cleanup)
			go h.handleIncomingWebsocketMessages(ctx, client, cleanup)
			h.handleOutgoingWebSocketMessages(ctx, client, cleanup)

		})(c)
	}
	return fiber.ErrUpgradeRequired
}

func (h *ChatHub) handlePingPong(ctx context.Context, client *Client, cleanup func()) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			client.mu.Lock()
			if client.closed {
				client.mu.Unlock()
				return
			}
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.mu.Unlock()
				log.Printf("Ping faild for user %s", client.UserID)
				cleanup()
				return
			}
			client.mu.Unlock()
		}
	}
}

func (h *ChatHub) handleIncomingWebsocketMessages(ctx context.Context, client *Client, cleanup func()) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg Message
			if err := client.Conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Websocket read error for user %s %v", client.UserID, err)
				}
				cleanup()
				return
			}

			kafkaMsg := &kafka.Message{
				FromUserID: msg.Sender,
				ToUserID:   msg.Recever,
				Content:    msg.Content,
				Timestamp:  time.Now().Unix(),
			}

			if err := h.SendMessageWithRetry(kafkaMsg, 3); err != nil {
				log.Printf("Error sending message form user %s: %v", client.UserID, err)
				client.Conn.WriteJSON(map[string]interface{}{
					"type":  "error",
					"error": "Faild to send meesage",
				})
			} else {
				log.Printf("message sent form %s to %s", msg.Sender, msg.Recever)
			}
		}
	}
}

func (h *ChatHub) handleOutgoingWebSocketMessages(ctx context.Context, client *Client, cleanup func()) {
	defer cleanup()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-client.Send:
			if !ok {
				return
			}

			client.mu.Lock()
			if client.closed {
				client.mu.Unlock()
				return
			}

			if err := client.Conn.WriteJSON(data); err != nil {
				client.mu.Unlock()
				log.Printf("Error writing to user %s : %v", client.UserID, err)
				return
			}
			client.mu.Unlock()
		}
	}
}
