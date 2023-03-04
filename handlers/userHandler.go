package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	dto "chat/DTO/userdto"
	ucontroller "chat/controllers/usersController"
	userDb "chat/entities/user"

	"golang.org/x/crypto/bcrypt"
)

var (
	errHandler     = errorHandler{}
	userService    = userDb.New()
	userController = ucontroller.New(userService)
)

// signing up new user
func signUp(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("couldn't parse data to CreateUserDTO err: %v", err)
			return
		}
		newUser := dto.CreateUserDTO{
			Name:     r.FormValue("name"),
			Mail:     r.FormValue("mail"),
			Password: r.FormValue("pass"),
		}

		// validating data
		if newUser.Name == "" {
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
			log.Fatal(err)
		}
		newUser.Password = string(hash)
		err = userController.CreateUser(newUser)
		if err != nil {
			errHandler.err = fmt.Errorf("server error")
			errHandler.responseForError(w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}

}

func loginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			Tmpl.ExecuteTemplate(w, "login.html", nil)
		}

		if r.Method == "POST" {
			user := dto.GetUserDTO{
				Mail:     r.FormValue("mail"),
				Password: r.FormValue("pass"),
			}
			// Validating user's data
			if user.Mail == "" {
				errHandler.err = fmt.Errorf("missing mail")
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}
			if user.Password == "" {
				errHandler.err = fmt.Errorf("password is too short, at least 8 elements required")
				errHandler.responseForError(w, http.StatusBadRequest)
				return
			}

			// Getting user from db
			dbUser, err := userController.GetUser(user)
			if err != nil {
				if err == sql.ErrNoRows {
					errHandler.err = fmt.Errorf("user is not found")
					errHandler.responseForError(w, http.StatusBadRequest)
					return
				} else {
					log.Fatal(err)
				}
			}

			// Comparing passwords
			err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
			if err != nil {
				errHandler.err = fmt.Errorf("invalid email or password")
				errHandler.responseForError(w, http.StatusUnauthorized)
				return
			}
			err = AuthUser(w, dbUser)
			if err != nil {
				errHandler.err = err
				errHandler.responseForError(w, 500)
			}
			http.Redirect(w, r, "/chat/", http.StatusSeeOther)
		}

	})

}
