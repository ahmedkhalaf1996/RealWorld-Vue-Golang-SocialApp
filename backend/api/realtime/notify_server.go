package realtime

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type Notification struct {
	ID        string    `json:"_id"`
	Details   string    `json:"details"`
	MainUID   string    `json:"mainuid"`
	TargetID  string    `json:"targetid"`
	IsReaded  bool      `json:"isreded"`
	CraetedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type User struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

type NotificationManager struct {
	connections map[string]*websocket.Conn
	lock        sync.RWMutex
}

var notificationManager *NotificationManager

func init() {
	notificationManager = &NotificationManager{
		connections: make(map[string]*websocket.Conn),
	}
}

func GetNotificationManager() *NotificationManager {
	return notificationManager
}

func (nm *NotificationManager) AddNotificationConnection(userID string, conn *websocket.Conn) {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	nm.connections[userID] = conn
	log.Printf("User %s connected to notitificaton server", userID)
}

func (nm *NotificationManager) RemoveNotificationConnection(userID string) {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	delete(nm.connections, userID)
	log.Printf("User %s disconnected from notitificaton server", userID)
}

func (nm *NotificationManager) SendNotificatonToUser(userID string, notifiaton Notification) error {
	nm.lock.RLock()
	conn, exists := nm.connections[userID]
	nm.lock.RUnlock()

	if !exists {
		log.Printf("User %s not connected for notification server", userID)
		return nil
	}

	err := conn.WriteJSON(notifiaton)
	if err != nil {
		log.Printf("Error sending notification to user %s : %v", userID, err)
		nm.RemoveNotificationConnection(userID)
		return err
	}

	log.Printf("notification sent to user %s : %s", userID, notifiaton.Details)
	return nil
}

func (nm *NotificationManager) GetConnectedUsers() []string {
	nm.lock.RLock()
	defer nm.lock.RUnlock()

	users := make([]string, 0, len(nm.connections))
	for userID := range nm.connections {
		users = append(users, userID)
	}
	return users
}
