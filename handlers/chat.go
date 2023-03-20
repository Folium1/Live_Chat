package handlers

import (
	auth "chat/handlers/middleware"
	logger "chat/logger"
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
	conn, _ := upgrader.Upgrade(w, r, nil)

	clients = append(clients, *conn)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		userMessage := []byte(user.Name + ": " + string(msg))
		for _, client := range clients {
			err := client.WriteMessage(msgType, userMessage)
			_,err = auth.GetToken(r)
			if err != nil {
				
			}
			if err != nil {
				logger.Error("Couldn't deliver msg", err, funcName)
			}
		}
	}
}
