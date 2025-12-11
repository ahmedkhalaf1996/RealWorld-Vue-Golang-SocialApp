package realtime

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaNotificationBridge struct {
	writer  *kafka.Writer
	reader  *kafka.Reader
	ctx     context.Context
	cancel  context.CancelFunc
	handler NotificationDeliveryHandler
}

type NotificationDeliveryHandler interface {
	DeliverToLocalClient(userID string, notification Notification)
}

func NewKafkaNotificationBridge(kafkaAddr, nodeID string) (*KafkaNotificationBridge, error) {
	ctx, cancel := context.WithCancel(context.Background())

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaAddr),
		Topic:                  "notifications",
		Balancer:               &kafka.Hash{},
		BatchTimeout:           10 * time.Millisecond,
		WriteTimeout:           10 * time.Second,
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

	bridge := &KafkaNotificationBridge{
		writer: writer,
		reader: reader,
		ctx:    ctx,
		cancel: cancel,
	}

	go bridge.consumeNotification()
	log.Println("kafka notificatin bridge initialized")
	return bridge, nil
}

func (k *KafkaNotificationBridge) SetDeliveryHandler(handler NotificationDeliveryHandler) {
	k.handler = handler
}

func (k *KafkaNotificationBridge) PublishNotification(targetUID string, notification Notification) error {
	dataBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := kafka.Message{
		Key:   []byte(targetUID),
		Value: dataBytes,
	}

	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("faild to publish notification to kafka : %v", err)
	}

	return nil
}

func (k *KafkaNotificationBridge) consumeNotification() {
	log.Println("kafka notificaton consumer started")

	for {
		select {
		case <-k.ctx.Done():
			log.Println("notification consumer stopeed")
			return
		default:
		}

		msg, err := k.reader.ReadMessage(k.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Printf("kafak read error %v, retrying in 1s", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var notification Notification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			log.Printf("Error unmarshaling notification %v", err)
			continue
		}

		if k.handler != nil {
			k.handler.DeliverToLocalClient(notification.MainUID, notification)
		}
	}
}

func (k *KafkaNotificationBridge) Close() {
	log.Println("Closing kafka notificaiton bridge...")
	k.cancel()

	if k.writer != nil {
		k.writer.Close()
	}
	if k.reader != nil {
		k.reader.Close()
	}
	log.Println("kafka notifcation bridge closed..")
}
