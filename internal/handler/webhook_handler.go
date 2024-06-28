package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/esc-chula/esc-discord-bot/config"
	"github.com/esc-chula/esc-discord-bot/internal/instance"
	"github.com/esc-chula/esc-discord-bot/pkg/utils"
)

type webhookHandler struct {
}

func NewWebhookHandler() *webhookHandler {
	return &webhookHandler{}
}

type WebhookPayload struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

func (h *webhookHandler) Received(w http.ResponseWriter, r *http.Request) {
	clientSecret := r.Header.Get("X-Bot-Secret")
	if clientSecret != config.GetConfig().Webhook.Secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload WebhookPayload

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	log.Print("Webhook received, updating user data...")

	usersData, err := utils.GetUsersData()
	if err != nil {
		http.Error(w, "Unable to get user data", http.StatusInternalServerError)
		return
	}
	instance.SetUsersInstance(usersData)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}
