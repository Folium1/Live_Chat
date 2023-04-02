package handlers

import (
	dto "chat/DTO/userdto"
	uController "chat/controllers/usersController"
	httpErr "chat/handlers/errorHandler"
	auth "chat/handlers/middleware"
	templParse "chat/handlers/templateParser"
	"chat/logger"

	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	errHandler     = httpErr.ErrorHandler{}
	userController = uController.New()
	tmpl           = templParse.TemplateParse()
)

// signUpHandler signing up new user
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()

	logger.InfoHttp(r.URL.Path, r.Method, funcName)
	// if auth.IsAuthenticated(w, r) {
	// 	http.Redirect(w, r, "/login/", http.StatusSeeOther)
	// }

	if r.Method == "GET" {
		err := tmpl.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			errHandler.Err = fmt.Errorf("server error")
			errHandler.ResponseForError(w, http.StatusInternalServerError)
			return
		}
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			logger.Error("couldn't parse form", err, funcName)
			return
		}
		// parsing data from the form
		newUser := dto.CreateUserDTO{
			Name:     r.FormValue("name"),
			Mail:     r.FormValue("mail"),
			Password: r.FormValue("pass"),
		}

		// validating data
		if newUser.Name == "" {
			errHandler.Err = fmt.Errorf("missing name")
			errHandler.ResponseForError(w, http.StatusBadRequest)
			return
		}
		if newUser.Mail == "" {
			errHandler.Err = fmt.Errorf("missing mail")
			errHandler.ResponseForError(w, http.StatusBadRequest)
			return
		}
		if len(newUser.Password) < 8 {
			errHandler.Err = fmt.Errorf("password is too short, at least 8 elements required")
			errHandler.ResponseForError(w, http.StatusBadRequest)
			return
		}

		// Hashing user's password
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		if err != nil {

			logger.Error("Couldn't hash user's password to check", err, funcName)
		}
		newUser.Password = string(hash)
		userId, err := userController.CreateUser(newUser)
		if err != nil {
			logger.Error("Couldn't create user", err, funcName)
			errHandler.Err = fmt.Errorf("server error")
			errHandler.ResponseForError(w, http.StatusInternalServerError)
			return
		}
		auth.AuthUser(w, userId)
		http.Redirect(w, r, "/chat/", http.StatusFound)
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)
	if auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/chat/", http.StatusSeeOther)
	}
	if r.Method == "GET" {

		err := tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {

			logger.Error("Couldn't parse template", err, funcName)
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
			errHandler.Err = fmt.Errorf("missing mail")

			logger.Error("No mail in user's input", nil, funcName)
			errHandler.ResponseForError(w, http.StatusBadRequest)
			return
		}
		if user.Password == "" {
			errHandler.Err = fmt.Errorf("password is too short, at least 8 elements required")

			logger.Error("No password in user's input", nil, funcName)
			errHandler.ResponseForError(w, http.StatusBadRequest)
			return
		}

		// Getting user from db
		dbUser, err := userController.GetUser(user)
		if err != nil {
			if err == sql.ErrNoRows {
				errHandler.Err = fmt.Errorf("user is not found")

				logger.Error("No user in db", err, funcName)
				errHandler.ResponseForError(w, http.StatusBadRequest)
				return
			} else {
				logger.Error("Couldn't create user", err, funcName)
			}
		}

		// Comparing passwords
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			errHandler.Err = fmt.Errorf("invalid email or password")
			logger.Error("Mail or password is invalid", err, funcName)
			errHandler.ResponseForError(w, http.StatusUnauthorized)
			return
		}
		err = auth.AuthUser(w, dbUser.Id)
		if err != nil {
			errHandler.Err = err
			logger.Error("Couldn't authorize user", err, funcName)
			errHandler.ResponseForError(w, 500)
		}
		http.Redirect(w, r, "/chat/", http.StatusFound)
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()
	logger.InfoHttp(r.URL.Path, r.Method, funcName)
	auth.LogOut(w, r)
	http.Redirect(w, r, "/login/", 302)
}
