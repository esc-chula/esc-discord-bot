package repository

import (
	"encoding/json"
	"fmt"

	"github.com/esc-chula/esc-discord-bot/config"
	"github.com/esc-chula/esc-discord-bot/pkg/http"
)

type UserRepository struct {
	client *http.CustomHTTPClient
}

func NewUserRepository() *UserRepository {
	cfg := config.GetConfig()

	client := http.NewCustomHTTPClient()
	client.SetBaseURL(fmt.Sprintf("%s/api/v2/tables/%s", cfg.Webhook.NocoDBAPIEndpoint, cfg.Webhook.NocoDBTableId))
	client.SetDefaultHeaders(map[string]string{
		"Accept":   "application/json",
		"xc-token": cfg.Webhook.NocoDBAPIToken,
	})
	client.SetDefaultContentType("application/json")

	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) GetUsersData() ([]map[string]interface{}, error) {
	resp, err := r.client.Get("/records?limit=1000")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	usersDataBytes, err := json.Marshal(data["list"])
	if err != nil {
		return nil, err
	}

	var usersData []map[string]interface{}

	err = json.Unmarshal(usersDataBytes, &usersData)
	if err != nil {
		return nil, err
	}

	for _, user := range usersData {
		if user["Discord ID"] != nil {
			user["Bot Status"] = "confirmed"
		} else {
			user["Bot Status"] = "unconfirmed"
		}
	}

	return usersData, nil
}

func (r *UserRepository) UpdateUserDiscordId(id interface{}, discordId string) error {
	payload := map[string]interface{}{
		"Id":         id,
		"Discord ID": discordId,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := r.client.Patch("/records", payloadBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}
