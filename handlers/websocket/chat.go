package socket

import (
	chatdto "chat/DTO/chatDTO"
	"chat/controllers/chatController"
	"chat/controllers/usersController"
	"chat/handlers/errorHandler"
	templParse "chat/handlers/templateParser"

	"chat/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	clients           []websocket.Conn
	messageController = chatController.New()
	tmpl              = templParse.TemplateParse()
	httpErr = errorHandler.ErrorHandler{}
	userController    = usersController.New()
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


// chatTemplate executes html template
func ChatTemplate(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)
	// loading history of previous massages
	prevMessages, err := messageController.GetAllMessages()
	if err != nil {
		logger.Error("Couldn't get all messages from db", err, funcName)
		http.Error(w, "Server error", 500)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", prevMessages)
}

func Chat(w http.ResponseWriter, r *http.Request) {
	// Log function name and HTTP request details
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)

	userId := r.Context().Value("userId").(string)
	// Get user information from DB by user ID
	user, err := userController.GetUserById(userId)
	if err != nil {
		errorHandler.HandleAuthError(w, r, funcName)
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

// addClient adds new client to clients slice
func addClient(conn *websocket.Conn) {
	clients = append(clients, *conn)
}

// createMessageAndBroadcast creates new message and broadcasts it to all clients
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
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)

	// Get the ID of the message to be deleted from the URL path
	id := strings.Split(r.URL.Path, "/")[2]
	if id == "" {
		logger.Error("No message id provided", fmt.Errorf("No message id provided"), funcName)
		httpErr.ResponseForError(w, 400)
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
