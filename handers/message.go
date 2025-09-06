package handers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/SomeSuperCoder/global-chat/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MessageHandler struct {
	Repo repository.MessageRepo
}

type GetResult struct {
	Messages   []models.Message `json:"messages"`
	TotalCount int64            `json:"total_count"`
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Prase data
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page == "" {
		http.Error(w, "No page number provided", http.StatusBadRequest)
		return
	}
	if limit == "" {
		http.Error(w, "No limit number provided", http.StatusBadRequest)
		return
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	limitNumber, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit number", http.StatusBadRequest)
		return
	}

	// Fetch messages
	messages, totalCount, err := h.Repo.FindPaged(r.Context(), int64(pageNumber), int64(limitNumber))
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Make a result
	result := &GetResult{
		Messages:   messages,
		TotalCount: totalCount,
	}
	resultString, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to from a proper response", http.StatusInternalServerError)
		return
	}

	// Send the result
	fmt.Fprintln(w, string(resultString))
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	userAuth := middleware.ExtractUserAuth(r)
	text := r.FormValue("text")

	if text == "" {
		http.Error(w, "The message text must not be empty!", http.StatusNotAcceptable)
		return
	}

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

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Get the message ID
	messageID := r.PathValue("id")
	if messageID == "" {
		panic("no message ID provided")
	}

	parsedMessageID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID provided", http.StatusBadRequest)
		return
	}
	// Delete the message
	h.Repo.DeleteMessage(r.Context(), parsedMessageID)
}

func (h *MessageHandler) UpdateMessageText(w http.ResponseWriter, r *http.Request) {
	// Get the message ID
	messageID := r.PathValue("id")
	if messageID == "" {
		panic("No message ID provided")
	}

	parsedMessageID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID provided", http.StatusBadRequest)
		return
	}

	// Get the message text
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	messageText := string(body)

	// Update the message
	err = h.Repo.UpdateMessageText(r.Context(), parsedMessageID, messageText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
