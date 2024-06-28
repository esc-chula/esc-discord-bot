package utils

import (
	"encoding/json"
	"fmt"

	"github.com/esc-chula/esc-discord-bot/config"
	"github.com/esc-chula/esc-discord-bot/pkg/http"
)

func GetUsersData() ([]map[string]interface{}, error) {
	client := http.NewCustomHTTPClient()
	client.SetBaseURL(config.GetConfig().Webhook.NocoDBAPIEndpoint)
	client.SetDefaultHeaders(map[string]string{
		"Accept":   "application/json",
		"xc-token": config.GetConfig().Webhook.NocoDBAPIToken,
	})

	resp, err := client.Get(fmt.Sprintf("/api/v2/tables/%s/records?limit=200", config.GetConfig().Webhook.NocoDBTableId))
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

	return usersData, nil
}
