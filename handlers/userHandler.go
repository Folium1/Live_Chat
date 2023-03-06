package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	dto "chat/DTO/userdto"
	ucontroller "chat/controllers/usersController"
	userDb "chat/entities/user"
	auth "chat/handlers/middleware"
	"chat/logger"

	"golang.org/x/crypto/bcrypt"
)

var (
	errHandler     = errorHandler{}
	userService    = userDb.New()
	userController = ucontroller.New(userService)
)

// signing up new user
func signUpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcName := logger.GetFuncName()
		logger.InfoHttp(r.URL.Path, r.Method, funcName)
		auth.IsAuthenticated(w, r)
		if r.Method == "GET" {
			err := Tmpl.ExecuteTemplate(w, "register.html", nil)
			if err != nil {
				errHandler.err = fmt.Errorf("server error")
				errHandler.responseForError(w, http.StatusInternalServerError)
				return
			}
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				funcName := logger.GetFuncName()
				logger.Error("couldn't parse form", err, funcName)
				return
			}
			newUser := dto.CreateUserDTO{
				Name:     r.FormValue("name"),
				Mail:     r.FormValue("mail"),
				Password: r.FormValue("pass"),
			}

			// validating data
			if newUser.Name < "" {
				errHandler.err = fmt.Errorf("missing name")
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}
			if newUser.Mail == "" {
				errHandler.err = fmt.Errorf("missing mail")
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}
			if len(newUser.Password) < 8 {
				errHandler.err = fmt.Errorf("password is too short, at least 8 elements required")
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}

			// Hashing user's password
			hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
			if err != nil {
				funcName := logger.GetFuncName()
				logger.Error("Coudn't hash user's password to check", err, funcName)
			}
			newUser.Password = string(hash)
			err = userController.CreateUser(newUser)
			if err != nil {
				errHandler.err = fmt.Errorf("server error")
				errHandler.responseForError(w, http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

func loginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcName := logger.GetFuncName()
		logger.InfoHttp(r.URL.Path, r.Method, funcName)
		auth.IsAuthenticated(w, r)
		if r.Method == "GET" {

			err := Tmpl.ExecuteTemplate(w, "login.html", nil)
			if err != nil {
				funcName := logger.GetFuncName()
				logger.Error("Couldn't parse tamplate", err, funcName)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		}

		if r.Method == "POST" {
			user := dto.GetUserDTO{
				Mail:     r.FormValue("mail"),
				Password: r.FormValue("pass"),
			}
			// Validating user's data
			if user.Mail == "" {
				errHandler.err = fmt.Errorf("missing mail")
				funcName := logger.GetFuncName()
				logger.Error("No mail in user's input", nil, funcName)
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}
			if user.Password == "" {
				errHandler.err = fmt.Errorf("password is too short, at least 8 elements required")
				funcName := logger.GetFuncName()
				logger.Error("No password in user's input", nil, funcName)
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}

			// Getting user from db
			dbUser, err := userController.GetUser(user)
			if err != nil {
				if err == sql.ErrNoRows {
					errHandler.err = fmt.Errorf("user is not found")
					funcName := logger.GetFuncName()
					logger.Error("No user in db", err, funcName)
					errHandler.responseForError(w, http.StatusBadRequest)
					return
				} else {
					funcName := logger.GetFuncName()
					logger.Error("Couldn't create user", err, funcName)
				}
			}

			// Comparing passwords
			err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
			if err != nil {
				errHandler.err = fmt.Errorf("invalid email or password")
				// l.Error(errHandler.err)
				funcName := logger.GetFuncName()
				logger.Error("Mail or password is invalid", err, funcName)
				errHandler.responseForError(w, http.StatusUnauthorized)
				return
			}
			err = auth.AuthUser(w, r, dbUser)
			if err != nil {
				errHandler.err = err
				funcName := logger.GetFuncName()
				logger.Error("Coudn't authenticate user", err, funcName)
				errHandler.responseForError(w, 500)
			}
			http.Redirect(w, r, "/chat/", http.StatusFound)
		}

	})

}
