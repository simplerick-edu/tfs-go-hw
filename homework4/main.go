package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Post("/login", Login)
	root.Post("/signup", Signup)

	r := chi.NewRouter()
	r.Use(Auth)
	r.Get("/messages", GetChatMessages)
	r.Post("/messages", PostChatMessages)
	r.Get("/{id}/messages", GetUserMessages)
	r.Post("/{id}/messages", PostUserMessages)

	root.Mount("/users/", r)

	log.Fatal(http.ListenAndServe(":5000", root))
}
