package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// WebSocketManager управляет подключениями пользователей
type WebSocketManager struct {
	clients    map[int]*websocket.Conn // Карта пользователей с их соединениями
	broadcast  chan Message
	register   chan Client
	unregister chan int
}

// Client представляет пользователя
type Client struct {
	ID     int
	Socket *websocket.Conn
}

// Message представляет сообщение между двумя пользователями
type Message struct {
	ID         string    `json:"_id"`
	SenderID   int       `json:"senderId"`
	ReceiverID int       `json:"receiverId"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Создание нового менеджера WebSocket
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[int]*websocket.Conn),
		broadcast:  make(chan Message),
		register:   make(chan Client),
		unregister: make(chan int),
	}
}

// Управление подключениями и сообщениями
func (manager *WebSocketManager) Run(db *sql.DB) {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client.ID] = client.Socket
			log.Printf("User %d connected", client.ID)
		case userID := <-manager.unregister:
			if conn, ok := manager.clients[userID]; ok {
				conn.Close()
				delete(manager.clients, userID)
				log.Printf("User %d disconnected", userID)
			}
		case message := <-manager.broadcast:
			// Сохранение сообщения в базу данных
			saveMessageToDB(db, message)

			// Отправка сообщения получателю
			if conn, ok := manager.clients[message.ReceiverID]; ok {
				err := conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error sending message to user %d: %v", message.ReceiverID, err)
					manager.unregister <- message.ReceiverID
				}
			}
		}
	}
}

// Обработчик подключения WebSocket
func (app *application) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	var clientData struct {
		UserID int `json:"userId"`
	}
	err = conn.ReadJSON(&clientData)
	if err != nil {
		log.Printf("Failed to read client data: %v", err)
		conn.Close()
		return
	}

	client := Client{
		ID:     clientData.UserID,
		Socket: conn,
	}
	app.wsManager.register <- client

	// Обработка входящих сообщений
	go app.handleWebSocketMessages(conn, clientData.UserID)
}

// Обработка сообщений от клиента
func (app *application) handleWebSocketMessages(conn *websocket.Conn, userID int) {
	defer func() {
		app.wsManager.unregister <- userID
		conn.Close()
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message:  %v", err)
			break
		}

		msg.ID = generateMessageID()
		msg.CreatedAt = time.Now()

		// Добавление сообщения в канал broadcast
		app.wsManager.broadcast <- msg
	}
}

// Сохранение сообщения в базу данных
func saveMessageToDB(db *sql.DB, msg Message) {
	_, err := db.Exec(`
		INSERT INTO messages (id, sender_id, receiver_id, text, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Text, msg.CreatedAt,
	)
	if err != nil {
		log.Printf("Error saving message to database: %v", err)
	}
}

// Генерация уникального ID для сообщения
func generateMessageID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}
	return sb.String()
}

func (app *application) WebSocketAIHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket AI upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		var userMsg struct {
			UserID int    `json:"userId"`
			Text   string `json:"text"`
		}
		if err := conn.ReadJSON(&userMsg); err != nil {
			log.Printf("Error reading AI message: %v", err)
			break
		}

		// Получаем ответ от AI
		aiResponse, err := getAIResponse(userMsg.Text)
		if err != nil {
			log.Printf("Error getting AI response: %v", err)
			break
		}

		// Готовим сообщение
		response := Message{
			ID:         generateMessageID(),
			SenderID:   0, // AI
			ReceiverID: userMsg.UserID,
			Text:       aiResponse,
			CreatedAt:  time.Now(),
		}

		// Отправляем клиенту
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error sending AI response: %v", err)
			break
		}

		// Сохраняем в базу
		saveMessageToDB(app.db, response)
	}
}

func getAIResponse(prompt string) (string, error) {
	// Замените на свой OpenAI ключ
	apiKey := os.Getenv("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/chat/completions"

	payload := strings.NewReader(`{
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "` + prompt + `"}]
	}`)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "AI не смог ответить", nil
	}
	return res.Choices[0].Message.Content, nil
}
