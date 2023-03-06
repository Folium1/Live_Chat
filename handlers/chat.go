package handlers

import (
	auth "chat/handlers/middleware"
	logger "chat/logger"
	"database/sql"
	"log"
	"net/http"
)

func chatHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcName := logger.GetFuncName()
		logger.InfoHttp(r.URL.Path, r.Method, funcName)
		// checking if token is present, if not - redirects to login page
		token, err := auth.GetToken(r)
		if err != nil {
			// l.Error(r.URL.Path + " " + r.Method + " " + err.Error())
			funcName := logger.GetFuncName()
			logger.Info("Couldn't get token", funcName)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		// getting user's id from token
		userId, err := auth.ValidateToken(token)
		if err != nil {
			funcName := logger.GetFuncName()
			logger.Error("Couldn't validate token", err, funcName)
			http.Error(w, "No user id", http.StatusFound)
		}
		user, err := userController.GetUserById(userId)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Redirect(w, r, "/", http.StatusNotFound)
				return
			} else {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		message := "Hello " + user.Name
		_, err = w.Write([]byte(message))
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	})
}
