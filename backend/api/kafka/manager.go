package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type KafkaManager struct {
	addr string
}

func NewKafkaManager(addr string) *KafkaManager {
	return &KafkaManager{addr: addr}
}

func (km *KafkaManager) EnsureTopics(topics []string) error {
	conn, err := kafka.Dial("tcp", km.addr)
	if err != nil {
		return fmt.Errorf("faild to connect to kafka:%v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("faild to get the controller: %v", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))
	if err != nil {
		return fmt.Errorf("faild to connect to controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{}
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		})
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

func (km *KafkaManager) HealthCheck() error {
	conn, err := kafka.Dial("tcp", km.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Brokers()
	return err
}

func WaitForKafka(addr string, timeout time.Duration) error {
	km := NewKafkaManager(addr)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := km.HealthCheck(); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("kafka not ready after %v", timeout)
}

func ReadAllPartitions(kafkaAddr, topic string, handler func(msg kafka.Message) error) (int, error) {
	log.Printf("Reading All partitions for topic:%s", topic)

	conn, err := kafka.Dial("tcp", kafkaAddr)
	if err != nil {
		return 0, fmt.Errorf("faild to connect to kafka : %v", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return 0, fmt.Errorf("faild to read partitions : %v", err)
	}

	log.Printf("Found %d partitons(s) for topic %s", len(partitions), topic)

	totalProcessed := 0

	for _, partition := range partitions {
		log.Printf("Reading partition %d for topic %s..", partition.ID, topic)

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{kafkaAddr},
			Topic:     topic,
			Partition: partition.ID,
			MinBytes:  1,
			MaxBytes:  10e6,
			MaxWait:   100 * time.Millisecond,
		})

		reader.SetOffset(0)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		partitionProcessed := 0

		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.DeadlineExceeded {
					log.Printf("partition %d: timeout reached after processing %d messages", partition.ID, partitionProcessed)
					break
				}
				// no msg in this partition
				break
			}

			if err := handler(msg); err != nil {
				cancel()
				reader.Close()
				return totalProcessed, err
			}

			totalProcessed++
			partitionProcessed++
		}

		cancel()
		reader.Close()

		log.Printf("Partition %d: processed %d messages", partition.ID, partitionProcessed)
	}

	log.Printf("Compleated reading all partitions for %s: %d total messages", topic, totalProcessed)
	return totalProcessed, nil
}
