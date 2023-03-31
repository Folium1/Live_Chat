package handlers

import (
	middleware "chat/handlers/middleware"
	logger "chat/logger"
	"html/template"
	"net/http"
)

var (
	Tmpl *template.Template
)

func StartServer() {
	var err error
	Tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't parse tamplate", err, funcName)
	}
	http.Handle("/ws/", middleware.AuthMiddleware(http.HandlerFunc(chat)))
	http.Handle("/chat/", middleware.AuthMiddleware(http.HandlerFunc(chatTemplate)))
	http.Handle("/delete/", middleware.AuthMiddleware(http.HandlerFunc(deleteMessage)))
	http.HandleFunc("/login/", loginHandler().ServeHTTP)
	http.HandleFunc("/sign-up/", signUpHandler().ServeHTTP)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't start the server", err, funcName)
	}
}
