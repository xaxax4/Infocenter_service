package services

import (
	"Infocenter/internal/models"
	"os"
	"strconv"
	"sync"
	"time"
)

type MiddleMan struct {
	Topics          map[string]*models.Topic
	idCounter       int
	mu              sync.RWMutex
	timeoutDuration time.Duration
}

func NewMiddleMan() *MiddleMan {
	timeoutStr := os.Getenv("CLIENT_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "30" // default 30 seconds
	}

	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeoutSeconds = 30 // default if error
	}

	return &MiddleMan{
		Topics:          make(map[string]*models.Topic),
		idCounter:       0,
		timeoutDuration: time.Duration(timeoutSeconds) * time.Second,
	}
}

func (m *MiddleMan) startTopic(name string) {
	m.mu.RLock()
	topic := m.Topics[name]
	m.mu.RUnlock()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-topic.Done():
			return

		case msg := <-topic.Broadcast:
			topic.Lock()
			for client := range topic.Clients {
				select {
				case client.Chat <- msg:
				default:
					m.removeClient(topic, client)
				}
			}
			topic.Unlock()

		case <-ticker.C:
			now := time.Now()
			topic.Lock()
			for client := range topic.Clients {
				if now.Sub(client.CreatedAt) > m.timeoutDuration {
					duration := now.Sub(client.CreatedAt) // .Round(time.Second)
					timeoutMsg := models.Message{
						ID:      -1,
						Content: duration.String(),
						Time:    now,
					}
					select {
					case client.Chat <- timeoutMsg:
					default:
					}
					m.removeClient(topic, client)
				}
			}
			topic.Unlock()
		}
	}
}

func (m *MiddleMan) removeClient(topic *models.Topic, client *models.Client) {
	_, exists := topic.Clients[client]
	if exists {
		close(client.Chat)
		delete(topic.Clients, client)
	}
}

func (m *MiddleMan) SendToClients(topicName string, str string) {
	m.mu.Lock()
	msg := models.Message{
		ID:      m.idCounter,
		Content: str,
		Time:    time.Now(),
	}
	m.idCounter++

	topic, exists := m.Topics[topicName]
	if !exists {
		topic = createTopic(topicName, m)
		go m.startTopic(topicName)
	}
	m.mu.Unlock()

	topic.Broadcast <- msg
}

func (m *MiddleMan) AddClient(topicName string) *models.Client {
	m.mu.Lock()
	topic, exists := m.Topics[topicName]
	if !exists {
		topic = createTopic(topicName, m)
		go m.startTopic(topicName)
	}
	m.mu.Unlock()

	client := &models.Client{
		Chat:      make(chan models.Message, 15),
		CreatedAt: time.Now(),
	}

	topic.Lock()
	topic.Clients[client] = struct{}{}
	topic.Unlock()

	return client
}

func (m *MiddleMan) RemoveClient(topicName string, client *models.Client) {
	m.mu.RLock()
	topic, exists := m.Topics[topicName]
	m.mu.RUnlock()

	if !exists {
		return
	}

	topic.Lock()
	defer topic.Unlock()
	_, exists = topic.Clients[client]
	if !exists {
		return
	}

	m.removeClient(topic, client)

	if len(topic.Clients) == 0 {
		topic.Close()
		m.mu.Lock()
		delete(m.Topics, topicName)
		m.mu.Unlock()
	}
}

func createTopic(topicName string, m *MiddleMan) *models.Topic {
	topic := models.NewTopic(topicName)
	m.Topics[topicName] = topic
	return topic
}
