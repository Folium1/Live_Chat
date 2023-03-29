package handlers

import (
	chatdto "chat/DTO/chatDTO"
	chatDb "chat/entities/chat"
	"chat/controllers/chatController"
	auth "chat/handlers/middleware"
	"chat/logger"
	"encoding/json"
	"strings"

	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var (
	clients           []websocket.Conn
	messageService    = chatDb.Message{}
	messageController = chatController.New(&messageService)
)

// chatTemplate executes html template
func chatTemplate(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	_, err := auth.GetToken(w, r)
	if err != nil {
		handleAuthError(w, r, funcName)
		return
	}
	// loading history of previous massages
	prevMessages, err := messageController.GetAllMessages()
	if err != nil {
		logger.Error("Couldn't get all messages from db", err, funcName)
		http.Error(w, "Server error", 500)
		return
	}
	Tmpl.ExecuteTemplate(w, "index.html", prevMessages)
}

func chat(w http.ResponseWriter, r *http.Request) {
	// Log function name and HTTP request details
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)

	// Get user token from cookies
	token, err := auth.GetToken(w, r)
	if err != nil {
		handleAuthError(w, r, funcName)
		return
	}

	// Validate user token and get user ID
	userId, err := auth.ValidateToken(token)
	if err != nil {
		handleAuthError(w, r, funcName)
		logger.Error("Couldn't validate token", err, funcName)
		return
	}

	// Get user information from DB by user ID
	user, err := userController.GetUserById(userId)
	if err != nil {
		handleAuthError(w, r, funcName)
		logger.Error(fmt.Sprintf("Couldn't get user by ID (%v)", userId), err, funcName)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Couldn't upgrade connection to WebSocket", err, funcName)
		return
	}
	defer conn.Close()

	// Add new client to clients slice
	addClient(conn)

	// Read messages from client
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		message := chatdto.CreateMessageDTO{
			Text:     string(msg),
			UserId:   userId,
			UserName: user.Name,
		}
		if err := createMessageAndBroadcast(message, msgType); err != nil {
			logger.Error("Couldn't create message", err, funcName)
			http.Error(w, "Server error", 500)
			return
		}
	}
}

func handleAuthError(w http.ResponseWriter, r *http.Request, funcName string) {
	auth.DeleteCookies(w)
	http.Redirect(w, r, "/login", 302)
}

func addClient(conn *websocket.Conn) {
	clients = append(clients, *conn)
}

func createMessageAndBroadcast(msg chatdto.CreateMessageDTO, msgType int) error {
	var err error
	msg.Id, err = messageController.CreateMsg(msg)
	if err != nil {
		return err
	}
	newMessageData, _ := json.Marshal(msg)
	for _, client := range clients {
		if err := client.WriteMessage(msgType, newMessageData); err != nil {
			logger.Error("Couldn't deliver message", err, logger.GetFuncName())
		}
	}
	return nil
}

// make a handle function to delete message by id from db
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)

	// Get the ID of the message to be deleted from the URL path
	id := strings.Split(r.URL.Path, "/")[2]
	if id == "" {
		logger.Error("No message id provided", fmt.Errorf("No message id provided"), funcName)
		http.Error(w, "Server error", 500)
		return
	}

	// Delete the message from the database
	err := messageController.DeleteMsg(id)
	if err != nil {
		logger.Error("Couldn't delete message", err, funcName)
		http.Error(w, "Server error", 500)
		return
	}
	for i, client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("Chat.Delete:"+id))
		if err != nil {
			// remove the disconnected client from the slice of clients
			clients = append(clients[:i], clients[i+1:]...)
			logger.Error("Couldn't deliver delete message to client", err, funcName)
		}
	}
	// Return a success message
	response := map[string]string{"message": "success"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Error("Couldn't marshal JSON response", err, funcName)
		http.Error(w, "Server error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
