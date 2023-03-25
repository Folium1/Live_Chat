package handlers

import (
	chatdto "chat/DTO/chatDTO"
	auth "chat/handlers/middleware"
	logger "chat/logger"
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
	clients []websocket.Conn
)

// chatTemplate executes html template
func chatTemplate(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	_, err := auth.GetToken(r)
	if err != nil {
		auth.DeleteCookies(w)
		http.Redirect(w, r, "/login", 302)
		return
	}
	// loading history of previous massages
	prevMessages, err := chController.GetAllMessages()
	if err != nil {
		logger.Error("Couldn't get all messages from db", err, funcName)
		http.Error(w, "Server error", 500)
		return
	}
	Tmpl.ExecuteTemplate(w, "index.html", prevMessages)
}

func chat(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)
	// getting token from user's cookies
	token, err := auth.GetToken(r)
	if err != nil {
		auth.DeleteCookies(w)
		http.Redirect(w, r, "/login/", 302)
		return
	}
	// validate token and getting user's id
	userId, err := auth.ValidateToken(token)
	if err != nil {
		auth.DeleteCookies(w)
		logger.Error("Couldn't validate token", err, funcName)
		http.Redirect(w, r, "/login/", 302)
		return
	}
	// getting user's information from db by his id
	user, err := userController.GetUserById(userId)
	if err != nil {
		auth.DeleteCookies(w)
		logger.Error(fmt.Sprintf("Couldn't get user by id(%v)", userId), err, funcName)
		http.Redirect(w, r, "/login/", 302)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Couldn't upgrade connection to WebSocket", err, funcName)
		return
	}
	defer conn.Close()

	clients = append(clients, *conn)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		message := chatdto.CreateMessageDTO{}
		message.Text = string(msg)
		message.UserId = userId
		message.UserName = user.Name
		chController.CreateMsg(message)
		for _, client := range clients {
			err := client.WriteMessage(msgType, []byte(message.UserName+": "+message.Text))
			_, err = auth.GetToken(r)
			if err != nil {

			}
			if err != nil {
				logger.Error("Couldn't deliver msg", err, funcName)
			}
		}
	}
}

// make a handle function to delete message by id from db
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id := parts[2]
	fmt.Println(id)
	if id == "" {
		logger.Error("Couldn't delete message", fmt.Errorf("No message id provided"), funcName)
		http.Error(w, "Server error", 500)
		return
	}
	err := chController.DeleteMsg(id)
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
