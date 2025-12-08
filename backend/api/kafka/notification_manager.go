package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Notification struct {
	UserID    string                 `json:"user_id"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	ID        string                 `json:"id"`
}

type NotificationHandler interface {
	DeliverNotification(notif *Notification)
}

type NotificationManager struct {
	kafkaWriter        *kafka.Writer
	notificationReader *kafka.Reader
	handler            NotificationHandler
	ctx                context.Context
	cancel             context.CancelFunc
}

func NewNotificationManager(kafkaAddr, nodeID string, handler NotificationHandler) (*NotificationManager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	km := NewKafkaManager(kafkaAddr)
	if err := km.EnsureTopics([]string{"notifications"}); err != nil {
		log.Printf("Warning: could not ensure notifications topic exists: %v", err)
	}

	time.Sleep(1 * time.Second)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaAddr),
		Balancer:               &kafka.Hash{},
		BatchTimeout:           10 * time.Millisecond,
		WriteTimeout:           11 * time.Second,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddr},
		Topic:          "notifications",
		GroupID:        "notification-group-" + nodeID,
		StartOffset:    kafka.LastOffset,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	nm := &NotificationManager{
		kafkaWriter:        writer,
		notificationReader: reader,
		handler:            handler,
		ctx:                ctx,
		cancel:             cancel,
	}

	go nm.listenToNotifications()
	log.Println("Notification manager initialized")
	return nm, nil
}

func (nm *NotificationManager) PublishNotification(notif *Notification) error {
	notif.Timestamp = time.Now().Unix()
	notif.ID = fmt.Sprintf("%s-%d", notif.UserID, notif.Timestamp)

	dataBytes, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	event := Event{
		Type: "notification",
		Data: json.RawMessage(dataBytes),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kafkaMsg := kafka.Message{
		Key:   []byte(notif.UserID),
		Value: eventBytes,
		Topic: "notifications",
	}
	if err := nm.kafkaWriter.WriteMessages(ctx, kafkaMsg); err != nil {
		return err
	}

	log.Printf("Notification published for user %s : %s", notif.UserID, notif.Title)
	return nil
}

func (nm *NotificationManager) listenToNotifications() {
	log.Println("Notification kafka listener stated")

	for {
		select {
		case <-nm.ctx.Done():
			log.Println("notification listern stopped")
			return
		default:
		}

		msg, err := nm.notificationReader.ReadMessage(nm.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Printf("kafka notifcation read error : %v, retrying in 1s", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshling notification event: %v", err)
			continue
		}

		if event.Type == "notification" {
			var notif Notification
			if err := json.Unmarshal(event.Data, &notif); err != nil {
				log.Printf("Error unmarshaling notification: %v", err)
				continue
			}

			log.Printf("Processing notifcation for user %s: %s", notif.UserID, notif.Title)
			nm.handler.DeliverNotification(&notif)
		}
	}
}

func (nm *NotificationManager) Close() {
	log.Println("Stopping notification manger..")
	nm.cancel()

	if nm.kafkaWriter != nil {
		nm.kafkaWriter.Close()
	}
	if nm.notificationReader != nil {
		nm.notificationReader.Close()
	}

	log.Println("notification manger stopped!")
}
