package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var Tmpl *template.Template

func chatHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid := r.Context().Value("userId").(int)
		message := fmt.Sprintf("Hello " + strconv.Itoa(userid))
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
	http.HandleFunc("/login/", AuthMiddleWare(loginHandler).ServeHTTP)
	http.HandleFunc("/sign-up/", signUp)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
