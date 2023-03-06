package handlers

import (
	"html/template"

	auth "chat/handlers/middleware"
	logger "chat/logger"

	"net/http"
	// "github.com/sirupsen/logrus"
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

	http.HandleFunc("/chat/", auth.AuthMiddleWare(chatHandler).ServeHTTP)
	http.HandleFunc("/", loginHandler().ServeHTTP)
	http.HandleFunc("/sign-up/", signUpHandler().ServeHTTP)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't start the server", err, funcName)
		// l.Fatal("Couldn't start the server, err: %v", err)
	}
}
