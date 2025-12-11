package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type NodeHeartbeat struct {
	NodeID    string   `json:"node_id"`
	TimeStamp int64    `json:"timestamp"`
	UserIDs   []string `json:"user_ids"`
}

type ChatHub interface {
	GetConnectedUserIDs() []string
	PublishUserStatus(userID, nodeID string, online bool) error
}

type HeartbeatManager struct {
	mu                 sync.RWMutex
	nodeID             string
	hub                ChatHub
	kafkaWriter        *kafka.Writer
	heartbeatReader    *kafka.Reader
	nodeLastSeen       map[string]int64
	processedDeadNodes map[string]bool
	nodeUsers          map[string][]string
	heartbeatInterval  time.Duration
	timeoutThershould  time.Duration
	ctx                context.Context
	cancel             context.CancelFunc
}

func NewHeartbeatManager(kafkaAddr, nodeID string, hub ChatHub) (*HeartbeatManager, error) {
	// func NewHeartbeatManager(kafkaAddr, nodeID string) (*HeartbeatManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	km := NewKafkaManager(kafkaAddr)
	if err := km.EnsureTopics([]string{"node-heartbeat"}); err != nil {
		log.Printf("Warning: could not ensure heartbeat topic exists: %v", err)
	}

	time.Sleep(1 * time.Second)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaAddr),
		Topic:                  "node-heartbeat",
		Balancer:               &kafka.Hash{},
		BatchTimeout:           10 * time.Millisecond,
		WriteTimeout:           10 * time.Second,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddr},
		Topic:          "node-heartbeat",
		GroupID:        "heartbeat-monitor-" + nodeID,
		StartOffset:    kafka.LastOffset,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	hm := &HeartbeatManager{
		nodeID:             nodeID,
		hub:                hub,
		kafkaWriter:        writer,
		heartbeatReader:    reader,
		nodeLastSeen:       make(map[string]int64),
		nodeUsers:          make(map[string][]string),
		processedDeadNodes: make(map[string]bool),
		heartbeatInterval:  5 * time.Second,
		timeoutThershould:  15 * time.Second,
		ctx:                ctx,
		cancel:             cancel,
	}

	///
	log.Printf("[%s] rebuilding node state form heartbeats..", nodeID)
	if err := hm.rebuildNodeStateFromAllPartitions(kafkaAddr); err != nil {
		log.Printf("[%s] node state rebuilding has issues: %v", nodeID, err)
	}

	// sendheatbeats
	go hm.sendheatbeats()
	// listenToHeartBeats
	go hm.listenToHeartBeats()
	// monitorNodeHealth
	go hm.monitorNodeHealth()

	// cleanupOldDeadNodes
	go hm.cleanupOldDeadNodes()
	log.Printf("Heatbeat amnager initazed for node %s", nodeID)
	return hm, nil
}

func (hm *HeartbeatManager) rebuildNodeStateFromAllPartitions(kafkaAddr string) error {
	log.Printf("[%s] reading all heatbeat messages from all partitons..", hm.nodeID)

	totalProcessed, err := ReadAllPartitions(kafkaAddr, "node-heartbeat", func(msg kafka.Message) error {
		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return nil
		}

		if event.Type == "node_heatbeat" {
			var heartbeat NodeHeartbeat
			if err := json.Unmarshal(event.Data, &heartbeat); err != nil {
				return nil
			}

			hm.mu.Lock()
			hm.nodeLastSeen[heartbeat.NodeID] = heartbeat.TimeStamp
			hm.nodeUsers[heartbeat.NodeID] = heartbeat.UserIDs
			hm.mu.Unlock()
		}

		return nil

	})

	if err != nil {
		return err
	}

	log.Printf("[%s] Node state rebuilding complated : %d totoal heartbeats", hm.nodeID, totalProcessed)
	hm.mu.RLock()
	for nodeID, lastSeen := range hm.nodeLastSeen {
		users := hm.nodeUsers[nodeID]
		log.Printf("[%s] Known node : %s (last need: %d , users: %d)", hm.nodeID, nodeID, lastSeen, len(users))
	}

	hm.mu.RUnlock()

	return nil
}

func (hm *HeartbeatManager) sendheatbeats() {
	ticker := time.NewTicker(hm.heartbeatInterval)
	defer ticker.Stop()

	log.Printf("Starting heatbeat sender for node %s", hm.nodeID)

	for {
		select {
		case <-hm.ctx.Done():
			log.Println("Heatbeat sender stopped")
			return
		case <-ticker.C:
			if err := hm.publishHeartbeat(); err != nil {
				log.Printf("Faild to publish heatbeat: %s", err)
			}
		}
	}
}

func (hm *HeartbeatManager) publishHeartbeat() error {
	userIDs := hm.hub.GetConnectedUserIDs()
	// userIDs := []string{""}

	heartbeat := NodeHeartbeat{
		NodeID:    hm.nodeID,
		TimeStamp: time.Now().Unix(),
		UserIDs:   userIDs,
	}

	data, err := json.Marshal(heartbeat)
	if err != nil {
		return err
	}

	event := Event{
		Type: "node_heartbeat",
		Data: json.RawMessage(data),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := kafka.Message{
		Key:   []byte(hm.nodeID),
		Value: eventBytes,
	}

	if err := hm.kafkaWriter.WriteMessages(ctx, msg); err != nil {
		return err
	}

	log.Printf("Heatbeat sent: node=%s, user=%d", hm.nodeID, len(userIDs))

	return nil
}

func (hm *HeartbeatManager) listenToHeartBeats() {
	log.Println("Heatbeat listener started")

	for {
		select {
		case <-hm.ctx.Done():
			log.Println("Heatbeat listener stoppped")
			return
		default:
		}

		msg, err := hm.heartbeatReader.ReadMessage(hm.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Printf("kafka heatbeat read error : %v, retring in 1s", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling heartbeat event :%v", err)
			continue
		}

		if event.Type == "node_heartbeat" {
			var heatbeat NodeHeartbeat
			if err := json.Unmarshal(event.Data, &heatbeat); err != nil {
				log.Printf("Error unmarshaling heatbeat :%v", err)
				continue
			}

			hm.handleHeartBeat(heatbeat)
		}

	}
}

func (hm *HeartbeatManager) handleHeartBeat(heartbeat NodeHeartbeat) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.nodeLastSeen[heartbeat.NodeID] = heartbeat.TimeStamp
	hm.nodeUsers[heartbeat.NodeID] = heartbeat.UserIDs

	if hm.processedDeadNodes[heartbeat.NodeID] {
		log.Printf("Node %s is Alive again! removing form dead list", heartbeat.NodeID)
		delete(hm.processedDeadNodes, heartbeat.NodeID)
	}

	log.Printf("Heatbeat received : node=%s, users=%d, timestamp=%d", heartbeat.NodeID, len(heartbeat.UserIDs), heartbeat.TimeStamp)
}

func (hm *HeartbeatManager) monitorNodeHealth() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Node health monitor started")

	for {
		select {
		case <-hm.ctx.Done():
			log.Println("Node health monitor stopped")
			return
		case <-ticker.C:
			hm.checkDeadNodes()
		}
	}
}

func (hm *HeartbeatManager) checkDeadNodes() {
	now := time.Now().Unix()
	timeoutSeconds := int64(hm.timeoutThershould.Seconds())

	hm.mu.Lock()
	deadNodes := []string{}

	for nodeID, lastSeen := range hm.nodeLastSeen {
		if nodeID == hm.nodeID {
			continue
		}

		timeSinceLastSeen := now - lastSeen
		if timeSinceLastSeen > timeoutSeconds {
			if !hm.processedDeadNodes[nodeID] {
				deadNodes = append(deadNodes, nodeID)
				log.Printf("Node %s is Dead (last seen %d seconds ago)", nodeID, timeSinceLastSeen)
			}
		}
	}

	hm.mu.Unlock()
	if len(deadNodes) > 0 {
		if hm.shouldHandleFailure() {
			for _, nodeID := range deadNodes {
				hm.HandleDeadNode(nodeID)
			}
		} else {
			log.Printf("Not Handling dead nodes %v - andther node is responsible", deadNodes)
		}
	}
}

func (hm *HeartbeatManager) shouldHandleFailure() bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	aliveNodes := []string{hm.nodeID}
	now := time.Now().Unix()
	timesoutSeconds := int64(hm.timeoutThershould.Seconds())

	for nodeID, lastSeen := range hm.nodeLastSeen {
		if nodeID == hm.nodeID {
			continue
		}
		timeSinceLastSeen := now - lastSeen
		if timeSinceLastSeen <= timesoutSeconds {
			aliveNodes = append(aliveNodes, nodeID)
		}
	}

	sort.Strings(aliveNodes)

	if len(aliveNodes) == 0 {
		return true
	}

	designatedNode := aliveNodes[0]
	isDesignated := designatedNode == hm.nodeID

	if isDesignated {
		log.Printf("This Node (%s) is Designated to handle dead nodes (alive noded: %v)", hm.nodeID, aliveNodes)
	}
	return isDesignated
}

func (hm *HeartbeatManager) HandleDeadNode(nodeID string) {
	hm.mu.Lock()
	userIDs := hm.nodeUsers[nodeID]
	hm.processedDeadNodes[nodeID] = true
	hm.mu.Unlock()

	if len(userIDs) == 0 {
		log.Printf("Dead Node %s had no users", nodeID)
		return
	}

	log.Printf("Handling Dead Node %s : Marking %d users as offline.", nodeID, len(userIDs))

	for _, userID := range userIDs {
		if err := hm.hub.PublishUserStatus(userID, nodeID, false); err != nil {
			log.Printf("Faild to publish offline status for user %s: %v", userID, err)
		} else {
			log.Printf("User %s marked offline (form dead node : %s)", userID, nodeID)
		}
	}
}

func (hm *HeartbeatManager) cleanupOldDeadNodes() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Dead Node Cleanup Service Started")

	for {
		select {
		case <-hm.ctx.Done():
			log.Println("Dead Node cleanup Servvice Stopped")
			return
		case <-ticker.C:
			now := time.Now().Unix()
			cleanupThereshould := int64(5 * 60)

			hm.mu.Lock()

			nodesToCleanup := []string{}
			for nodeID, lastSeen := range hm.nodeLastSeen {
				if nodeID == hm.nodeID {
					continue
				}

				timeSinceLastSeen := now - lastSeen
				if timeSinceLastSeen > cleanupThereshould && hm.processedDeadNodes[nodeID] {
					nodesToCleanup = append(nodesToCleanup, nodeID)
				}
			}

			// remove old dead nodes form memory
			for _, nodeID := range nodesToCleanup {
				log.Printf("forgetting old dead node : %s (last seedn %d) seconds agao", nodeID, now-hm.nodeLastSeen[nodeID])
				delete(hm.nodeLastSeen, nodeID)
				delete(hm.nodeUsers, nodeID)
				delete(hm.processedDeadNodes, nodeID)
			}

			hm.mu.Unlock()

			if len(nodesToCleanup) > 0 {
				log.Printf("Clean up %d old dead nodes", len(nodesToCleanup))
			}
		}
	}

}

// stop func for gracful shoultdonw
func (hm *HeartbeatManager) Stop() {
	log.Println("Stopping heatbeat manager..")
	hm.cancel()

	if hm.kafkaWriter != nil {
		hm.kafkaWriter.Close()
	}
	if hm.heartbeatReader != nil {
		hm.heartbeatReader.Close()
	}
	log.Println("Heatbet manager stopped")
}
