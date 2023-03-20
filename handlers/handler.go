package handlers

import (
	"fmt"
	"html/template"

	auth "chat/handlers/middleware"
	logger "chat/logger"

	"net/http"
	// "github.com/sirupsen/logrus"
)

var (
	Tmpl *template.Template
)

func mockChat() http.Handler {
	funcName := logger.GetFuncName()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoHttp(r.URL.Path, r.Method, funcName)
		token, err := auth.GetToken(r)
		if err != nil {
			logger.Error("Couldn't get token from cookies", err, funcName)
			http.Redirect(w, r, "/login", http.StatusNotFound)
			return
		}
		userId, err := auth.ValidateToken(token)
		if err != nil {
			logger.Error("Couldn't validate token", err, funcName)
			http.Redirect(w, r, "/login", http.StatusNotFound)
			return
		}
		user, err := userController.GetUserById(userId)
		if err != nil {
			logger.Error(fmt.Sprintf("coudn't get user by id(%v)", userId), err, funcName)
			http.Error(w, "couldn't get user data", 404)
			return
		}
		msg := fmt.Sprintf("Hello,%v!!!", user.Name)
		_, err = w.Write([]byte(msg))
		if err != nil {
			logger.Error("couldn't write message to user", err, funcName)
		}
	})

}

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
