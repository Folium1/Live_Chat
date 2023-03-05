package handlers

import (
	"html/template"
	"log"
	"net/http"
)

var Tmpl *template.Template

func chatHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checking if token is present, if not - redirects to login page
		token, err := getToken(r)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		// getting user's id from token
		userId, err := validateToken(token)
		if err != nil {
			http.Error(w, "No user id", http.StatusFound)
		}

		message := "Hello " + userId
		w.Write([]byte(message))
	})
}

func StartServer() {
	var err error
	Tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/chat/", AuthMiddleWare(chatHandler).ServeHTTP)
	http.HandleFunc("/", loginHandler().ServeHTTP)
	http.HandleFunc("/sign-up/", AuthMiddleWare(signUpHandler).ServeHTTP)

	err = http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal(err)
	}
}
