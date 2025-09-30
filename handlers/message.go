package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SomeSuperCoder/global-chat/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var validate = validator.New()

type MessageHandler struct {
	Repo repository.MessageRepo
}

type MessageResponse struct {
	Messages   []models.Message `json:"messages"`
	TotalCount int64            `json:"total_count"`
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Prase
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

	// Do work
	messages, totalCount, err := h.Repo.FindPaged(r.Context(), int64(pageNumber), int64(limitNumber))
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Respond
	result := &MessageResponse{
		Messages:   messages,
		TotalCount: totalCount,
	}
	resultString, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to from a proper response", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(resultString))
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Parse
	var request struct {
		Text string `json:"text" bson:"text,omitempty" validate:"required,min=1,max=500"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// Validate
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON validation failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// Do work
	err = h.Repo.CreateMessage(r.Context(), models.Message{
		Author:   userAuth.UserID,
		Text:     request.Text,
		CratedAt: time.Now(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond
	fmt.Fprintf(w, "Message successfully created")
}

func (h *MessageHandler) UpdateMessageText(w http.ResponseWriter, r *http.Request) {
	// Parse path params
	messageID := r.PathValue("id")

	parsedMessageID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID provided", http.StatusBadRequest)
		return
	}
	// Parse body
	var request struct {
		Text string `json:"text" bson:"text,omitempty" validate:"omitempty,min=1,max=500"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// Validate
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON validation failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// Do work
	err = h.Repo.UpdateMessage(r.Context(), parsedMessageID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Parse
	messageID := r.PathValue("id")

	parsedMessageID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID provided", http.StatusBadRequest)
		return
	}

	// Delete the message
	h.Repo.DeleteMessage(r.Context(), parsedMessageID)
}
