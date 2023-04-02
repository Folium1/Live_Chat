package handlers

import (
	middleware "chat/handlers/middleware"
	templateparser "chat/handlers/templateParser"
	userHandler "chat/handlers/user"
	socket "chat/handlers/websocket"
	logger "chat/logger"

	"net/http"
)

var (
	tmpl = templateparser.TemplateParse()
)

func StartServer() {
	var err error
	http.Handle("/ws/", middleware.AuthMiddleware(http.HandlerFunc(socket.Chat)))
	http.Handle("/chat/", middleware.AuthMiddleware(http.HandlerFunc(socket.ChatTemplate)))
	http.Handle("/delete/", middleware.AuthMiddleware(http.HandlerFunc(socket.DeleteMessage)))
	http.HandleFunc("/login/", userHandler.LoginHandler)
	http.HandleFunc("/sign-up/", userHandler.SignUpHandler)
	http.HandleFunc("/logout/", userHandler.LogoutHandler)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't start the server", err, funcName)
	}
}
