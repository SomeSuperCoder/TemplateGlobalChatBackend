package handers

import (
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
)

type MessageHandler struct {
	Repo repository.MessageRepo
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	userAuth := middleware.ExtractUserAuth(r)
	text := r.FormValue("text")

	err := h.Repo.CreateMessage(r.Context(), models.Message{
		Author:   userAuth.UserID,
		Text:     text,
		CratedAt: time.Now(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
