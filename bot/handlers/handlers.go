package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type BotService interface {
	Start() error
	Stop()
}

type BotHandler struct {
	service BotService
}

func New(service BotService) *BotHandler {
	return &BotHandler{service}
}

func (b *BotHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/start", b.start)
		r.Post("/stop", b.stop)
		//r.Post("/restart_with_new_settings", b.changeSettings)
	})
	return r
}

func (b *BotHandler) start(w http.ResponseWriter, r *http.Request) {
	err := b.service.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (b *BotHandler) stop(w http.ResponseWriter, r *http.Request) {
	b.service.Stop()
	w.WriteHeader(http.StatusOK)
}
