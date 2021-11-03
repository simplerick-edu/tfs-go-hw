package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	root := chi.NewRouter()
	service := NewChatService()
	root.Use(middleware.Logger)
	root.Post("/login", service.Login)
	root.Post("/signup", service.Signup)

	r := chi.NewRouter()
	r.Use(service.Auth)
	r.Get("/messages", service.GetChatMessages)
	r.Post("/messages", service.PostChatMessages)
	r.Get("/{id}/messages", service.GetUserMessages)
	r.Post("/{id}/messages", service.PostUserMessages)

	root.Mount("/users/", r)

	log.Fatal(http.ListenAndServe(":5000", root))
}
