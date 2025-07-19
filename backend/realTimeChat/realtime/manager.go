package realtime

import (
	"log"
	"realTimeChat/servegrpc"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Sender  string `json:"sender"`
	Recever string `json:"recever"`
	Content string `json:"content"`
}

type ConnectionManager struct {
	connections    map[string]*websocket.Conn
	onlineFriends  map[string][]string
	getUserFriends func(string) <-chan []string
	lock           sync.Mutex
}

func NewConnectionManager(getUserFriends func(string) <-chan []string) *ConnectionManager {
	if getUserFriends != nil {
		return &ConnectionManager{
			connections:    make(map[string]*websocket.Conn),
			onlineFriends:  make(map[string][]string),
			getUserFriends: getUserFriends,
		}
	}
	return nil
}

func (cm *ConnectionManager) AddConnection(userID string, conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.connections[userID] = conn
	cm.onlineFriends[userID] = []string{}

	// Notifiy exsiting users about the new online friend
	for friendID := range cm.onlineFriends {
		if friendID != userID && cm.isFriend(userID, friendID) {
			cm.onlineFriends[friendID] = append(cm.onlineFriends[friendID], userID)
			err := cm.connections[friendID].WriteJSON(map[string]interface{}{
				"onlineFriends": cm.onlineFriends[friendID],
			})
			if err != nil {
				log.Printf("Error notifiying %s about %s : %v", friendID, userID, err)
				return
			}
		}
	}

	// update the online friends list for the new user
	go func() {
		for friends := range cm.getUserFriends(userID) {
			if friends == nil {
				continue
			}

			for _, friendID := range friends {
				if cm.connections[friendID] != nil {
					cm.onlineFriends[userID] = append(cm.onlineFriends[userID], friendID)
					err := cm.connections[userID].WriteJSON(map[string]interface{}{
						"onlineFriends": cm.onlineFriends[userID],
					})
					if err != nil {
						log.Printf("Error notifiying %s about %s : %v", userID, friendID, err)
						return
					}
				}
			}
		}
	}()
}

func (cm *ConnectionManager) RemoveConnection(userID string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connections, userID)
	delete(cm.onlineFriends, userID)
	for friendID := range cm.onlineFriends {
		for i, id := range cm.onlineFriends[friendID] {
			if id == userID {
				cm.onlineFriends[friendID] = append(cm.onlineFriends[friendID][:i], cm.onlineFriends[friendID][i+1:]...)
				err := cm.connections[friendID].WriteJSON(map[string]interface{}{
					"onlineFriends": cm.onlineFriends[friendID],
				})
				if err != nil {
					log.Printf("Error notifiying %s about %s : %v", friendID, userID, err)
				}
				break
			}
		}
	}
}

func (cm *ConnectionManager) SendToReceiver(msg Message) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if conn, ok := cm.connections[msg.Recever]; ok {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error Sening message to %s : %v", msg.Recever, err)
		}
		// Save Message To DB Via GRPC
		err = servegrpc.SendMessageClient(msg.Sender, msg.Recever, msg.Content)
		if err != nil {
			log.Fatalf("error Saving message to gRPC %s :%v", msg.Recever, err)
		}
	} else {
		log.Printf("Recever %s not found ", msg.Recever)
	}
}

// helper func is Friend
func (cm *ConnectionManager) isFriend(userID, friendID string) bool {
	friends := <-cm.getUserFriends(userID)
	if friends == nil {
		return false
	}
	for _, f := range friends {
		if f == friendID {
			return true
		}
	}
	return false
}
