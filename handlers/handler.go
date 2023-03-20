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
	http.HandleFunc("/ws/", chat)
	http.HandleFunc("/chat/", func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.GetToken(r)
		if err != nil {
			auth.DeleteCookies(w)
			http.Redirect(w, r, "/login", 302)
			return
		}
		Tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	http.HandleFunc("/login", loginHandler().ServeHTTP)
	http.HandleFunc("/sign-up/", signUpHandler().ServeHTTP)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't start the server", err, funcName)
	}
}
